package manager

import (
	"dfs/protogen/node"
	"dfs/rpc/client"
	"dfs/utils"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

// StorageNode is the node that stores the chunks
type StorageNode struct {
	ID            string
	Address       string
	LastHeartbeat time.Time
	// chunks is the chunks that this node has
	// if the value is true, the chunk is available
	// otherwise, the chunk is about to be uploaded
	chunks        map[uint64]bool
	FreeSpace     uint64 // free space, gathered from node
	TotalSpace    uint64 // total space, gathered from node
	RequiredSpace uint64 // required space for upload new files

	// management
	wg  sync.WaitGroup
	m   *Manager
	cli *client.Client
}

// Manager is the manager of the storage nodes, and chunks
type Manager struct {
	mu sync.Mutex

	nodes         map[string]*StorageNode
	chunks        map[uint64]*ChunkInfo
	files         map[string]FileMetadata
	candiateFiles map[string]FileMetadata

	// constant
	chunkSize    uint64
	replicaCount uint64

	wg    sync.WaitGroup
	close chan struct{}
}

// closeTimer is a helper function to close a timer
func closeTimer(t *time.Timer) {
	if !t.Stop() {
		<-t.C
	}
}

// NewManager creates a new manager
func NewManager(chunkSize uint64, replicaCount uint64) *Manager {
	m := &Manager{
		nodes:         make(map[string]*StorageNode),
		chunks:        make(map[uint64]*ChunkInfo),
		files:         make(map[string]FileMetadata),
		candiateFiles: make(map[string]FileMetadata),
		chunkSize:     chunkSize,
		replicaCount:  replicaCount,
		close:         make(chan struct{}),
	}
	m.wg.Add(3)
	go m.CheckCandiateFiles()
	go m.CheckNodeHealth()
	go m.MaintainChunkDuplicate()
	return m
}

// Close closes the manager
func (m *Manager) Close() {
	close(m.close)
	m.wg.Wait()
}

// preRemoveDeadNode the hook to remove dead node
// the lock is already acquired
func (m *Manager) preRemoveDeadNode(node *StorageNode) {
	// go through the chunks and remove the node from the chunk
	for chunk := range node.chunks {
		m.chunks[chunk].removeReplica(node.ID)
	}
}

// removeChunk removes the chunk from the node
func (sn *StorageNode) removeChunk(chunkId uint64) bool {
	chunk, ok := sn.m.chunks[chunkId]
	if !ok {
		panic(fmt.Sprintf("chunk %d not found", chunkId))
	}
	if v, ok := sn.chunks[chunkId]; ok {
		if !v {
			sn.RequiredSpace -= chunk.Size
		}
		delete(sn.chunks, chunkId)
		return v
	}
	return false
}

// runAsyncTask runs the passed function asynchronously
func (sn *StorageNode) runAsyncTask(f func(sn *StorageNode) error) {
	sn.wg.Add(1)
	go func() {
		defer sn.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Println("recovered from panic:", r)
			}
		}()
		if err := f(sn); err != nil {
			log.Println("error in async task:", err)
		}
	}()
}

// setChunkStored sets the givne chunk as stored
func (sn *StorageNode) setChunkStored(chunkId uint64) {
	v, ok := sn.chunks[chunkId]
	if ok && !v {
		sn.RequiredSpace -= sn.m.chunks[chunkId].Size
	}
	sn.chunks[chunkId] = true
}

// weight returns the storage node weight,
// which is the ratio of free space to total space
func (sn *StorageNode) weight() float32 {
	free := float32(sn.FreeSpace)
	if free <= 0 || sn.TotalSpace <= 0 {
		return 0
	}
	return free / float32(sn.TotalSpace)
}

// addChunk adds the chunk to the node
func (sn *StorageNode) addChunk(chunkId uint64, dummy bool) {
	if v, ok := sn.chunks[chunkId]; ok {
		if v == dummy {
			panic("should not happen")
		}
		// pass
		return
	}
	sn.chunks[chunkId] = !dummy
	if dummy {
		sn.RequiredSpace += sn.m.chunks[chunkId].Size
	}
}

// selectReplicaNode selects 'count' replicas for the given chunk
func (m *Manager) selectReplicaNode(ci *ChunkInfo, count int) ([]string, error) {
	type nodeWeight struct {
		node   *StorageNode
		weight float32
	}
	var nodeWeights []nodeWeight
	replicas := ci.Replicas.ToSet()
	for _, node := range m.nodes {
		if _, ok := replicas[node.ID]; ok {
			continue
		}
		w := node.weight()
		if node.FreeSpace-node.RequiredSpace < ci.Size || w == 0 {
			continue
		}
		nodeWeights = append(nodeWeights, nodeWeight{
			node:   node,
			weight: w,
		})
	}
	if len(nodeWeights) < count {
		return nil, errors.New("no enough node available")
	}
	sort.Slice(nodeWeights, func(i, j int) bool {
		return nodeWeights[i].weight > nodeWeights[j].weight
	})
	nodes := make([]string, 0, count)
	for i := range count {
		nodes = append(nodes, nodeWeights[i].node.ID)
	}
	return nodes, nil
}

// createChunk creates a new chunk
func (m *Manager) createChunk(chunkSize uint64, fm *FileMetadata) ChunkInfo {
	// generate a chunk id
	for {
		rid := rand.Uint64()
		if _, ok := m.chunks[rid]; !ok {
			return ChunkInfo{
				ChunkID:  rid,
				Replicas: utils.NewUnorderedList[string](),
				Size:     chunkSize,
				fm:       fm,
			}
		}
	}
}

// availableSpace returns the available space
func (m *Manager) availableSpace() uint64 {
	var total uint64
	for _, sn := range m.nodes {
		total += sn.FreeSpace - sn.RequiredSpace
	}
	return total
}

// GetFile returns the file metadata for the given key
func (m *Manager) GetFile(key string) (*FileMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fm, ok := m.files[key]
	if !ok {
		return nil, fmt.Errorf("file %s not found", key)
	}
	fm = fm.Clone()
	return &fm, nil
}

// DeleteFile deletes the file metadata for the given key
func (m *Manager) DeleteFile(key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	isCandiate := false
	fm, ok := m.files[key]
	if !ok {
		fm, ok = m.candiateFiles[key]
		if !ok {
			return false, fmt.Errorf("file %s not found", key)
		}
		isCandiate = true
	}
	for _, chunk := range fm.Chunks {
		m.preRemoveChunk(chunk.ChunkID)
		delete(m.chunks, chunk.ChunkID)
	}
	if isCandiate {
		delete(m.candiateFiles, key)
	} else {
		delete(m.files, key)
	}
	return isCandiate, nil
}

// ListFiles lists the files with the given prefix
func (m *Manager) ListFiles(prefix string) ([]FileMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ans := make([]FileMetadata, 0)
	for key, fm := range m.files {
		if strings.HasPrefix(key, prefix) {
			ans = append(ans, fm)
		}
	}
	return ans, nil
}

// FinishStoreFile finishes the store file
func (m *Manager) FinishStoreFile(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	fm, ok := m.candiateFiles[key]
	if !ok {
		return fmt.Errorf("file %s not found", key)
	}
	// check if the file is finished
	for _, chunk := range fm.Chunks {
		if !m.nodes[chunk.Replicas.GetUnderlyingList()[0]].chunks[chunk.ChunkID] {
			return fmt.Errorf("chunk %d is not finished", chunk.ChunkID)
		}
	}
	delete(m.candiateFiles, key)
	m.files[key] = fm
	return nil
}

// StoreFile creates the file metadata for a file with given key to be stored
func (m *Manager) StoreFile(key string, size uint64) (*FileMetadata, error) {
	chunks := (size + m.chunkSize - 1) / m.chunkSize
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.files[key]; ok {
		return nil, fmt.Errorf("file %s already exists", key)
	}
	if _, ok := m.candiateFiles[key]; ok {
		return nil, fmt.Errorf("file %s is about to be uploaded, you can delete the tmp file by 'delete' command", key)
	}
	if space := m.availableSpace(); chunks > space {
		return nil, fmt.Errorf("manager: not enough chunks available, need %d, but only %d available", chunks, space)
	}
	// we can store the file
	fm := FileMetadata{
		Name:      key,
		Size:      size,
		CreatedAt: time.Now(),
	}
	for i := uint64(0); i < chunks; i++ {
		var chunkSize uint64
		if i == chunks-1 {
			chunkSize = size - (m.chunkSize * i)
		} else {
			chunkSize = m.chunkSize
		}
		ci := m.createChunk(chunkSize, &fm)
		nodes, err := m.selectReplicaNode(&ci, int(m.replicaCount))
		if err != nil {
			m.removeDeadFile(fm)
			return nil, err
		}
		m.chunks[ci.ChunkID] = &ci
		for _, nodeId := range nodes {
			ci.addReplica(nodeId)
			m.nodes[nodeId].addChunk(ci.ChunkID, true)
		}
		fm.Chunks = append(fm.Chunks, ci)
	}
	m.candiateFiles[key] = fm
	// clone the file metadata
	fm = fm.Clone()
	return &fm, nil
}

// generateNodeId generates the node id
// it must be the node address for now
func (m *Manager) generateNodeId(address string) string {
	return address
}

// AddNode adds a new node to the manager
func (m *Manager) AddNode(address string) (s StorageNode, err error) {
	nodeId := m.generateNodeId(address)
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.nodes[nodeId]; ok {
		err = fmt.Errorf("node %s already exists", nodeId)
		return
	}
	var conn net.Conn
	conn, err = net.Dial("tcp", address)
	if err != nil {
		err = fmt.Errorf("failed to connect to node %s: %w", address, err)
		return
	}
	s = StorageNode{
		ID:            nodeId,
		Address:       address,
		LastHeartbeat: time.Now(),
		chunks:        make(map[uint64]bool),
		cli:           client.NewClient(conn),
	}
	s.m = m
	m.nodes[nodeId] = &s
	return
}

// UpdateHeartbeat updates the heartbeat of the node and returns
// the chunks that should be removed.
func (m *Manager) UpdateHeartbeat(nodeId string, metadata *node.Metadata) ([]uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	node, ok := m.nodes[nodeId]
	if !ok {
		return nil, fmt.Errorf("node %s not found", nodeId)
	}

	chunks := make(map[uint64]struct{}, len(metadata.Chunks))

	// do something for the chunk fields
	removeChunks := make([]uint64, 0)
	for _, chunk := range metadata.Chunks {
		chunks[chunk] = struct{}{}
		if v, ok := node.chunks[chunk]; !ok {
			// the chunk should be removed
			removeChunks = append(removeChunks, chunk)
		} else {
			if v {
				// the chunks is stored
				// pass
			} else {
				node.setChunkStored(chunk)
			}
		}
	}
	// remove unstored chunks
	for chunkId := range node.chunks {
		if _, ok := chunks[chunkId]; !ok {
			m.chunks[chunkId].removeReplica(nodeId)
			node.removeChunk(chunkId)
		}
	}

	node.LastHeartbeat = time.Now()
	node.TotalSpace = metadata.Totalspace
	node.FreeSpace = metadata.Freespace

	return removeChunks, nil
}

// UpdateChunkUploaded updates the chunk metadata and returns the error.
func (m *Manager) UpdateChunkUploaded(chunkId uint64, hash string, size uint64, replicas []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.chunks[chunkId]
	if !ok {
		return fmt.Errorf("chunk %d not found", chunkId)
	}
	if c.Size != size {
		return fmt.Errorf("chunk %d size mismatch", chunkId)
	}
	has := c.Replicas.ToSet()
	for _, node := range replicas {
		if nd, ok := m.nodes[node]; ok {
			nd.setChunkStored(chunkId)
			if _, ok := has[node]; !ok {
				c.Replicas.AddNoDuplicate(node)
			}
		}
	}
	c.Hash = hash
	return nil
}

// preRemoveChunk removes the chunk from the node and send the request
// to the node to remove the chunk if necessary.
func (m *Manager) preRemoveChunk(chunkId uint64) {
	if c, ok := m.chunks[chunkId]; ok {
		for _, nodeId := range c.Replicas.GetUnderlyingList() {
			nd := m.nodes[nodeId]
			if exist := nd.removeChunk(chunkId); !exist {
				continue
			}
			nd.runAsyncTask(func(sn *StorageNode) error {
				_, err := sn.cli.SendRequest(&node.DeleteChunk{Id: chunkId})
				return err
			})
		}
	}
}

// removeDeadFile removes the file from the manager and the node.
func (m *Manager) removeDeadFile(fm FileMetadata) {
	for _, chunk := range fm.Chunks {
		m.preRemoveChunk(chunk.ChunkID)
		delete(m.chunks, chunk.ChunkID)
	}
	delete(m.candiateFiles, fm.Name)
}

// GetNode returns the node by nodeId.
func (m *Manager) GetNode(nodeId string) (*StorageNode, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if node, exist := m.nodes[nodeId]; exist {
		return node, nil
	} else {
		return nil, errors.New("node not found")
	}
}

// ListNodes returns the list of nodes.
func (m *Manager) ListNodes() ([]*StorageNode, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	nodes := make([]*StorageNode, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// CheckCandiateFiles checks the candiate files and remove the dead files,
// which means the file is not uploaded within 30 minutes.
func (m *Manager) CheckCandiateFiles() {
	defer m.wg.Done()
	const duration = time.Minute
	for t := time.NewTimer(duration); true; t.Reset(duration) {
		select {
		case <-m.close:
			closeTimer(t)
			return
		case <-t.C:
		}
		var deadFiles []string
		exp := time.Now().Add(-30 * time.Minute)
		m.mu.Lock()
		for key, file := range m.candiateFiles {
			if file.CreatedAt.Before(exp) {
				deadFiles = append(deadFiles, key)
			}
		}
		for _, key := range deadFiles {
			m.removeDeadFile(m.candiateFiles[key])
		}
		m.mu.Unlock()
	}
}

// CheckNodeHealth checks the node health and remove the dead nodes.
func (m *Manager) CheckNodeHealth() {
	defer m.wg.Done()
	const duration = 10 * time.Second
	for t := time.NewTimer(duration); true; t.Reset(duration) {
		select {
		case <-m.close:
			closeTimer(t)
			return
		case <-t.C:
		}
		var deadNodes []string
		now := time.Now()
		m.mu.Lock()
		for _, node := range m.nodes {
			if now.Sub(node.LastHeartbeat) > 30*time.Second {
				deadNodes = append(deadNodes, node.ID)
			}
		}
		if len(deadNodes) > 0 {
			fmt.Printf("[controller] %d nodes are dead: %v\n", len(deadNodes), deadNodes)
		}
		for _, nodeID := range deadNodes {
			m.preRemoveDeadNode(m.nodes[nodeID])
			sn := m.nodes[nodeID]
			go func(sn *StorageNode) {
				sn.wg.Wait()
				sn.cli.Close()
			}(sn)
			delete(m.nodes, nodeID)
		}
		m.mu.Unlock()
	}
}

// MaintainChunkDuplicate maintains the chunk replicas.
func (m *Manager) MaintainChunkDuplicate() {
	defer m.wg.Done()
	const duration = 20 * time.Second
	for t := time.NewTimer(duration); true; t.Reset(duration) {
		select {
		case <-m.close:
			closeTimer(t)
			return
		case <-t.C:
		}
		var chunksNeedDuplicate [][]uint64
		m.mu.Lock()
		counts := 0
		// only check the file that is already stored on server
		for _, fm := range m.files {
			for _, chunk := range fm.Chunks {
				if chunk.Replicas.Len() < int(m.replicaCount) {
					needs := int(m.replicaCount) - chunk.Replicas.Len()
					for needs >= len(chunksNeedDuplicate) {
						chunksNeedDuplicate = append(chunksNeedDuplicate, []uint64{})
					}
					chunksNeedDuplicate[needs] = append(chunksNeedDuplicate[needs], chunk.ChunkID)
					counts++
				}
			}
		}
		if len(chunksNeedDuplicate) == 0 {
			m.mu.Unlock()
			continue
		}
		fmt.Printf("[controller] detected %d chunks need duplicate\n", counts)
		type info struct {
			chunkId  uint64
			replicas []string
		}
		// select the node that can store the chunk and gather the
		// tasks for each node, which can prevent from sending too many
		// RPC requests to the same node
		sends := make(map[string][]info)

	outer:
		for i, chunks := range chunksNeedDuplicate {
			if i == 0 {
				continue
			}
			for _, chunk := range chunks {
				c := m.chunks[chunk]
				replicas, err := m.selectReplicaNode(c, i)
				if err != nil {
					fmt.Printf("[controller] faild to select replica node for %d: %v\n", chunk, err)
					break outer
				}
				servers := c.Replicas.GetUnderlyingList()
				if len(servers) == 0 {
					panic(fmt.Sprintf("[controller] faild to find a server that contains chunkd %d\n", chunk))
				}
				server := servers[0]
				sends[server] = append(sends[server], info{
					chunkId:  chunk,
					replicas: replicas,
				})
			}
		}
		for server, chunks := range sends {
			fmt.Printf("[controller] create send %d chunk task(s) to node %s\n", len(chunks), server)
			// this will be okay as the chunk is a copy after go 1.21
			// use async task to avoid deadlock and other issues
			m.nodes[server].runAsyncTask(func(sn *StorageNode) error {
				cs := make([]*node.SendChunk, 0, len(chunks))
				for _, c := range chunks {
					cs = append(cs, &node.SendChunk{
						Id:       c.chunkId,
						Replicas: c.replicas,
					})
				}
				_, err := sn.cli.SendRequest(&node.SendChunks{
					Chunks: cs,
				})
				return err
			})
		}
		m.mu.Unlock()
	}
}
