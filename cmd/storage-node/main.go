package main

import (
	"dfs/protogen/common"
	"dfs/protogen/node"
	rpcClient "dfs/rpc/client"
	"dfs/rpc/server"
	"dfs/utils"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/proto"
)

var (
	// "/bigdata/students/" + $(whoami)
	StoragePathPrefix string
)

var (
	errChunkHashMismatch = errors.New("chunk hash mismatch")
	errChunkNotFound     = errors.New("chunk not found")
)

func getUsername() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.Username
}

func generateStoragePath(listenAddr string) string {
	addrReplacer := regexp.MustCompile(`[:\[\]]+`)
	folder := addrReplacer.ReplaceAllString(listenAddr, "_")
	folder = fmt.Sprintf("storage-node-%s", folder)
	return path.Join(StoragePathPrefix, folder)
}

type ChunkInfo struct {
	ChunkID     uint64
	Hash        string
	Size        uint64
	lastChecked int64
}

type mux struct {
	id            string
	storagePath   string
	cli           *rpcClient.Client
	listenAddress string

	chunks         map[uint64]*ChunkInfo
	candiateChunks map[uint64]*ChunkInfo

	mutex sync.Mutex
}

func main() {
	controllerAddr := flag.String("controllerAddr", "localhost:8081", "Controller address")
	listenAddr := flag.String("listenAddr", "localhost:9000", "Listen address (this part )")
	flag.StringVar(&StoragePathPrefix, "storePath", "/bigdata/students/"+getUsername(), "Storage path prefix")
	flag.Parse()
	s := server.NewServer()
	s.SetLogPrefix("[storage-node]")

	conn, err := net.Dial("tcp", *controllerAddr)
	if err != nil {
		panic(err)
	}

	client := rpcClient.NewClient(conn)

	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		panic(err)
	}

	storagePath := generateStoragePath(*listenAddr)

	// remove all files
	if err := os.RemoveAll(storagePath); err != nil {
		panic(err)
	}
	// create the directory for storage
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		panic(err)
	}

	fmt.Printf("[storage-node] start at %s, storage path: %s\n", *listenAddr, storagePath)

	m := mux{
		id:             *listenAddr,
		storagePath:    storagePath,
		cli:            client,
		listenAddress:  *listenAddr,
		chunks:         make(map[uint64]*ChunkInfo),
		candiateChunks: make(map[uint64]*ChunkInfo),
	}

	// register the RPC handlers
	s.RegisterByMessage(&node.PutChunk{}, m.putChunkRequest)
	s.RegisterByMessage(&node.GetChunk{}, m.getChunkRequest)
	s.RegisterByMessage(&node.SendChunk{}, m.sendChunkRequest)
	s.RegisterByMessage(&node.SendChunks{}, m.sendChunksRequest)
	s.RegisterByMessage(&node.DeleteChunk{}, m.deleteChunkRequest)

	// serve the RPC
	go s.Serve(listener)

	// register to controller
	if err = m.register(); err != nil {
		panic(err)
	}
	// check chunk health
	go m.checkChunks()
	// heartbeat
	m.sendHeartbeat()
}

// updateLastChecked updates the last checked time of the chunk
func (c *ChunkInfo) updateLastChecked(t int64) {
	if c.lastChecked < t {
		c.lastChecked = t
	}
}

// checkChunks checks the health of the chunks every 20 minutes
func (m *mux) checkChunks() {
	const duration = 1 * time.Minute
	for timer := time.NewTimer(duration); true; timer.Reset(duration) {
		<-timer.C
		chunksNeedToBeCheck := make([]*ChunkInfo, 0)
		now := time.Now()
		// check chunks every 20 minutes
		expiredTimeStamp := now.Add(-time.Minute * 20).Unix()
		m.mutex.Lock()
		for _, chunk := range m.chunks {
			if chunk.lastChecked <= expiredTimeStamp {
				chunksNeedToBeCheck = append(chunksNeedToBeCheck, chunk)
			}
		}
		m.mutex.Unlock()

		chunksNeedToBeRemoved := make([]*ChunkInfo, 0)
		chunksChecked := make([]*ChunkInfo, 0)
		var hash string
		for _, chunk := range chunksNeedToBeCheck {
			fi, err := os.Open(m.getFilename(chunk.ChunkID))
			if err == nil {
				hash, err = utils.HashReader(fi)
				fi.Close()
				if err != nil {
					continue
				}
				if hash != chunk.Hash {
					err = errChunkHashMismatch
				}
			}
			if err != nil {
				fmt.Printf("[storage-node] the chunk %d is broken, remove it from storage, reason: %v\n", chunk.ChunkID, err)
				chunksNeedToBeRemoved = append(chunksNeedToBeRemoved, chunk)
			} else {
				chunksChecked = append(chunksChecked, chunk)
			}
		}
		timestamp := now.Unix()
		m.mutex.Lock()
		for _, chunk := range chunksNeedToBeRemoved {
			m.removeChunk(chunk.ChunkID, true, nil)
		}
		for _, chunk := range chunksChecked {
			chunk.updateLastChecked(timestamp)
		}
		m.mutex.Unlock()
	}
}

// sendHeartbeat sends heartbeat to controller every 10 seconds
// it will contains the chunks the node owns and the disk space usage.
//
// After receiving the response, the controller will tell the node to
// remove the chunks that are not needed. And the node will remove
// those chunks.
func (m *mux) sendHeartbeat() {
	const duration = 10 * time.Second
	for timer := time.NewTimer(duration); true; timer.Reset(duration) {
		<-timer.C
		metadata, err := m.gatherMetadata()
		if err != nil {
			fmt.Printf("update heartbeat error: %v\n", err)
		} else {
			resp, err := m.cli.SendRequest(&node.Heartbeat{
				Metadata: metadata,
			})
			if err != nil {
				fmt.Printf("update heartbeat error: %v\n", err)
			} else {
				fmt.Println("heartbeat sent")
				respMsg := resp.(*node.HeartbeatResponse)
				reason := errors.New("removed by controller")
				m.mutex.Lock()
				for _, chunkId := range respMsg.RemoveChunks {
					m.removeChunk(chunkId, true, reason)
				}
				m.mutex.Unlock()
			}
		}
	}
}

// findChunkById finds the chunk by id
// it will find the chunk in the stored chunks (the `chunks` field),
// and the candiate chunks (the `candiateChunks` field)
func (m *mux) findChunkById(chunkId uint64) (c *ChunkInfo, isCandiate bool, ok bool) {
	c, ok = m.chunks[chunkId]
	if ok {
		return c, false, true
	}
	if c, ok = m.candiateChunks[chunkId]; ok {
		return c, true, true
	}
	return
}

// getFilename returns the storage path of the chunk
func (m *mux) getFilename(chunkId uint64) string {
	return filepath.Join(m.storagePath, strconv.FormatUint(chunkId, 10))
}

// putChunkRequest handles the PutChunk RPC request
func (m *mux) putChunkRequest(ctx server.Context, req proto.Message) (resp proto.Message, err error) {
	reqMsg := req.(*node.PutChunk)
	// generate hash
	var hash string
	hash, err = utils.HashBytes(reqMsg.Chunk.Data)
	if err != nil {
		return nil, err
	}
	if reqMsg.Chunk.Size != nil {
		if len(reqMsg.Chunk.Data) != int(*reqMsg.Chunk.Size) {
			return nil, errors.New("chunk size mismatch")
		}
	} else {
		reqMsg.Chunk.Size = proto.Uint64(uint64(len(reqMsg.Chunk.Data)))
	}
	m.mutex.Lock()
	if chunk, _, ok := m.findChunkById(reqMsg.Chunk.Id); ok {
		m.mutex.Unlock()
		if chunk.Hash != hash {
			return nil, errChunkHashMismatch
		}
		return &common.EmptyResponse{}, nil
	}
	reqMsg.Chunk.Hash = &hash
	// add the chunk to the candidate map
	ci := ChunkInfo{
		ChunkID:     reqMsg.Chunk.Id,
		Hash:        hash,
		Size:        *reqMsg.Chunk.Size,
		lastChecked: time.Now().Unix(),
	}
	m.candiateChunks[reqMsg.Chunk.Id] = &ci
	m.mutex.Unlock()
	defer func() {
		m.mutex.Lock()
		if err == nil {
			m.chunks[reqMsg.Chunk.Id] = m.candiateChunks[reqMsg.Chunk.Id]
		}
		delete(m.candiateChunks, reqMsg.Chunk.Id)
		m.mutex.Unlock()
	}()
	filename := m.getFilename(reqMsg.Chunk.Id)
	err = os.WriteFile(filename, reqMsg.Chunk.Data, 0644)
	if err != nil {
		return nil, err
	}
	// send to controller and other nodes if needed
	if reqMsg.Replicas[0] == m.listenAddress {
		// send to other node
		otherReplicas := reqMsg.Replicas[1:]
		errs := make([]chan error, len(otherReplicas))
		for i, replica := range otherReplicas {
			errs[i] = make(chan error, 1)
			go func(replica string, req *node.PutChunk) {
				defer close(errs[i])
				if _, err := utils.SendSingleRequest(replica, req); err != nil {
					errs[i] <- err
				} else {
					errs[i] <- nil
				}
			}(replica, reqMsg)
		}
		joinErrs := make([]error, 0, len(otherReplicas))
		for _, err := range errs {
			joinErrs = append(joinErrs, <-err)
		}
		replicas := make([]string, 0, len(reqMsg.Replicas))
		replicas = append(replicas, m.id)
		for i := range joinErrs {
			if joinErrs[i] == nil {
				replicas = append(replicas, otherReplicas[i])
			}
		}
		// else report to controller
		// only the primary node will report to controller
		_, err = m.cli.SendRequest(&node.ChunkUploaded{
			Id:       ci.ChunkID,
			Size:     ci.Size,
			Replicas: replicas,
			Hash:     hash,
		})
		if err != nil {
			return nil, err
		}
	}

	return &common.EmptyResponse{}, nil
}

// onceFunc will only call the passed function once
// when the returned function is called.
func onceFunc(f func()) func() {
	done := false
	return func() {
		if !done {
			f()
			done = true
		}
	}
}

// sendChunk will send the chunk to the nodes specified in the replicas,
// and unlock the mutex.
func (m *mux) sendChunk(c *ChunkInfo, replicas []string) ([]string, error) {
	unlock := onceFunc(m.mutex.Unlock)
	defer unlock()
	filename := m.getFilename(c.ChunkID)
	data, err := os.ReadFile(filename)
	if err != nil {
		// the chunk is broken, remove it
		m.removeChunk(c.ChunkID, true, err)
		return nil, err
	}
	// do not prevent other tasks
	unlock()
	// check checksum
	hash, err := utils.HashBytes(data)
	if err != nil {
		return nil, err
	}
	if hash != c.Hash {
		// the chunk is broken, remove it
		m.removeChunk(c.ChunkID, false, errChunkHashMismatch)
		return nil, errChunkHashMismatch
	}
	// send the chunk to other nodes
	successNode := make([]string, 0, len(replicas))
	cp := &node.PutChunk{
		Chunk: &node.Chunk{
			Id:   c.ChunkID,
			Data: data,
			Hash: &c.Hash,
			Size: &c.Size,
		},
	}
	for _, node := range replicas {
		if node == m.id {
			continue
		}
		// only contains the destination node so the destination
		// node will report the acquisition of the chunk
		cp.Replicas = []string{node}
		if _, err := utils.SendSingleRequest(node, cp); err != nil {
			log.Printf("send chunk to %s failed, %v\n", node, err)
		} else {
			successNode = append(successNode, node)
		}
	}
	return successNode, nil
}

// sendChunkRequest handles the SendChunk RPC request
func (m *mux) sendChunkRequest(ctx server.Context, req proto.Message) (proto.Message, error) {
	reqMsg := req.(*node.SendChunk)
	m.mutex.Lock()
	c, ok := m.chunks[reqMsg.Id]
	if !ok {
		m.mutex.Unlock()
		return nil, errChunkNotFound
	}
	// the following function call will unlock the mutex
	successNode, err := m.sendChunk(c, reqMsg.Replicas)
	if err != nil {
		return nil, err
	}
	return &node.SendChunkResponse{
		Replicas: successNode,
	}, nil
}

// sendChunksRequest handles the SendChunks RPC request
func (m *mux) sendChunksRequest(ctx server.Context, req proto.Message) (proto.Message, error) {
	reqMsg := req.(*node.SendChunks)
	resp := &node.SendChunksResponse{}
	for _, chunk := range reqMsg.Chunks {
		m.mutex.Lock()
		c, ok := m.chunks[chunk.Id]
		if !ok {
			resp.Chunks = append(resp.Chunks, &node.SendChunkResponse{})
			m.mutex.Unlock()
			continue
		}
		// the following function call will unlock the mutex
		successNode, err := m.sendChunk(c, chunk.Replicas)
		if err != nil {
			fmt.Printf("[storage-node] failed to send chunk %d: %v\n", chunk.Id, err)
		}
		resp.Chunks = append(resp.Chunks, &node.SendChunkResponse{
			Replicas: successNode,
		})
	}
	return resp, nil
}

// deleteChunkRequest handles the DeleteChunk RPC request
func (m *mux) deleteChunkRequest(ctx server.Context, req proto.Message) (proto.Message, error) {
	deleteChunk := req.(*node.DeleteChunk)
	m.removeChunk(deleteChunk.Id, false, nil)
	return &common.EmptyResponse{}, nil
}

// removeChunk removes the chunk from the storage
func (m *mux) removeChunk(chunkId uint64, locked bool, reason error) {
	if !locked {
		m.mutex.Lock()
		defer m.mutex.Unlock()
	}
	if _, ok := m.chunks[chunkId]; !ok {
		return
	}
	if reason != nil {
		fmt.Printf("[storage-node] the chunk %d is broken, remove it from storage, reason: %v\n", chunkId, reason)
	}
	delete(m.chunks, chunkId)
	os.Remove(m.getFilename(chunkId))
}

// getChunkRequest handles the GetChunk RPC request
func (m *mux) getChunkRequest(ctx server.Context, req proto.Message) (proto.Message, error) {
	reqMsg := req.(*node.GetChunk)
	m.mutex.Lock()
	c, ok := m.chunks[reqMsg.Id]
	if !ok {
		m.mutex.Unlock()
		return nil, errChunkNotFound
	}
	// read file within the mutex
	// actually, we can open the file and unlock the mutex
	// and then read the file outside the mutex
	// This should be okay on Linux, but not on Windows
	bytes, err := os.ReadFile(m.getFilename(reqMsg.Id))
	if err != nil {
		m.removeChunk(reqMsg.Id, true, err)
		m.mutex.Unlock()
		return nil, err
	}
	m.mutex.Unlock()
	// calculate checksum
	hash, err := utils.HashBytes(bytes)
	if err != nil {
		return nil, err
	}
	if hash != c.Hash {
		m.removeChunk(reqMsg.Id, false, errChunkHashMismatch)
		return nil, errChunkHashMismatch
	}
	return &node.Chunk{
		Id:   reqMsg.Id,
		Data: bytes,
		Hash: &hash,
		Size: &c.Size,
	}, nil
}

// getFreespace returns the free space of the storage
// it only support Unix-like system
func (m *mux) getFreespace() (uint64, error) {
	var stat unix.Statfs_t
	if err := unix.Statfs(m.storagePath, &stat); err != nil {
		return 0, err
	}
	return stat.Bavail * uint64(stat.Bsize), nil
}

// gatherMetadata gathers the metadata of the storage node
func (m *mux) gatherMetadata() (*node.Metadata, error) {
	freespace, err := m.getFreespace()
	if err != nil {
		return nil, err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	chunkIds := make([]uint64, 0, len(m.chunks))
	usedSpace := uint64(0)
	for _, chunk := range m.chunks {
		usedSpace += chunk.Size
		chunkIds = append(chunkIds, chunk.ChunkID)
	}
	metadata := &node.Metadata{
		Id:         m.id,
		Totalspace: freespace + usedSpace,
		Freespace:  freespace,
		Chunks:     chunkIds,
	}
	return metadata, nil
}

// register registers the storage node to the controller
func (m *mux) register() error {
	metadata, err := m.gatherMetadata()
	if err != nil {
		return err
	}
	req := &node.Register{
		ListenAddr: m.listenAddress,
		Metadata:   metadata,
	}
	_, err = m.cli.SendRequest(req)
	return err
}
