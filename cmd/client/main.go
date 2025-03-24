package main

import (
	"dfs/cmd/client/pfile"
	"dfs/cmd/client/taskmgr"
	"dfs/protogen/client"
	"dfs/protogen/common"
	"dfs/protogen/node"
	rpcClient "dfs/rpc/client"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

type Context struct {
	cli *rpcClient.Client
}

func printUsage() {
	fmt.Printf(`Usage of %s: %s [OPTIONS] <COMMAND> [ARGS]
Options:
  -csa <address>             Controller server address.
  -h                         Print this help message.
Commands:
  upload <filename> <key>    Upload the file to remote file system with specified key.
  download <key> <filename>  Download the file from remote file system with specified key.
  delete <key>               Delete the file from remote file system with specified key.
  list [prefix]              List all files in remote file system with optional prefix.
  node [nodeId]              Show the node info, if the nodeId is not specified, it will show all nodes.
`, os.Args[0], os.Args[0])
}

func checkCommand(commands []string) error {
	if len(commands) == 0 {
		return errors.New("no command specified")
	}
	if commands[0] == "upload" || commands[0] == "download" {
		if len(commands) != 3 {
			return fmt.Errorf("%s command requires 2 arguments", commands[0])
		}
	} else if commands[0] == "delete" {
		if len(commands) != 2 {
			return errors.New("delete command requires 1 argument")
		}
	} else if commands[0] == "list" {
		if len(commands) > 2 {
			return errors.New("list command requires at most 1 argument")
		}
	} else if commands[0] == "node" {
		if len(commands) > 2 {
			return errors.New("node command requires at most 1 argument")
		}
	} else {
		return errors.New("unknown command")
	}
	return nil
}

func main() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	controllerServerAddr := flagSet.String("csa", "localhost:8080", "Controller server address")
	flagSet.Usage = printUsage
	flagSet.Parse(os.Args[1:])
	err := checkCommand(flagSet.Args())
	if err != nil {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", *controllerServerAddr)
	if err != nil {
		panic(err)
	}
	cli := rpcClient.NewClient(conn)
	defer cli.Close()
	ctx := Context{
		cli: cli,
	}
	var operation string
	args := flagSet.Args()
	switch args[0] {
	case "upload":
		err = ctx.StoreFile(args[1], args[2])
		operation = "upload file"
	case "download":
		err = ctx.DownloadFile(args[1], args[2])
		operation = "download file"
	case "delete":
		err = ctx.DeleteFile(args[1])
		operation = "delete file"
	case "list":
		prefix := ""
		if len(args) == 2 {
			prefix = args[1]
		}
		err = ctx.ListFiles(prefix)
		operation = "list files"
	case "node":
		var nodeId *string
		if len(args) == 2 {
			nodeId = &args[1]
		}
		err = ctx.ListNodes(nodeId)
	default:
	}
	if err != nil {
		fmt.Printf("%s failed: %v\n", operation, err)
		os.Exit(1)
	}
}

type taskMgrContext struct {
	file *pfile.PFile
	cli  *rpcClient.Client
	mgr  *taskmgr.TaskManager
}

type fileRange struct {
	start   int64
	end     int64
	id      uint64
	servers []string
}

// handleUpload handles the upload task
// it will try to upload the file to the servers in the fileRange.servers
func (ctx *Context) handleUpload(c any, taskCtx any) error {
	tmc := c.(*taskMgrContext)
	fileRange := taskCtx.(*fileRange)
	data, err := tmc.file.ReadPart(fileRange.start, fileRange.end)
	if err != nil {
		return err
	}
	_, err = tmc.cli.SendRequest(&node.PutChunk{
		Chunk: &node.Chunk{
			Id:   fileRange.id,
			Data: data,
		},
		Replicas: fileRange.servers,
	})
	if err != nil {
		// try another server
		if len(fileRange.servers) <= 1 {
			return err
		}
		fileRange.servers = fileRange.servers[1:]
		tmc.mgr.AddTaskContext(fileRange.servers[0], fileRange)
		return nil
	}
	return nil
}

// gatherAllServers gathers all servers from the chunks
func gatherAllServers(chunks []*common.ChunkInfo) []string {
	servers := make(map[string]struct{})
	for _, chunk := range chunks {
		for _, server := range chunk.Replicas {
			servers[server] = struct{}{}
		}
	}
	res := make([]string, 0, len(servers))
	for server := range servers {
		res = append(res, server)
	}
	return res
}

func closeClientsFromTaskMgr(taskMgr *taskmgr.TaskManager, servers []string) {
	for _, server := range servers {
		ctx := taskMgr.GetContext(server)
		if ctx != nil {
			ctx.(*taskMgrContext).cli.Close()
		}
	}
}

// StoreFile uploads the file to the remote file system
func (ctx *Context) StoreFile(filename string, key string) error {
	// open the file to upload
	pfile, err := pfile.NewPFile(filename, pfile.ReadPart)
	if err != nil {
		return err
	}
	defer pfile.Close()
	// get file size
	fi, err := pfile.Stat()
	if err != nil {
		return err
	}
	fileSize := fi.Size()
	resp, err := ctx.cli.SendRequest(&client.StoreRequest{
		Key:  key,
		Size: uint64(fileSize),
	})
	if err != nil {
		return err
	}
	respMsg := resp.(*client.StoreResponse)
	allServers := gatherAllServers(respMsg.Chunks)

	taskMgr := taskmgr.NewTaskManager(ctx.handleUpload)
	defer closeClientsFromTaskMgr(taskMgr, allServers)
	for _, server := range allServers {
		conn, err := net.Dial("tcp", server)
		if err != nil {
			return err
		}
		taskMgr.AddContext(server, &taskMgrContext{
			file: pfile,
			cli:  rpcClient.NewClient(conn),
			mgr:  taskMgr,
		})
	}
	rangeStart := uint64(0)
	for _, chunk := range respMsg.Chunks {
		end := rangeStart + chunk.Size
		fileRange := fileRange{
			start:   int64(rangeStart),
			end:     int64(end),
			id:      chunk.Id,
			servers: chunk.Replicas,
		}
		rangeStart = end
		taskMgr.AddTaskContext(chunk.Replicas[0], &fileRange)
	}
	err = taskMgr.Run()
	if err == nil {
		_, err = ctx.cli.SendRequest(&client.StoreRequest{
			Key:      key,
			Size:     uint64(fileSize),
			Finished: true,
		})
	}
	if err == nil {
		fmt.Printf("file %s uploaded successfully\n", key)
	}
	return err
}

// handleDownload handles the download task
func (ctx *Context) handleDownload(c any, taskCtx any) error {
	tmc := c.(*taskMgrContext)
	fileRange := taskCtx.(*fileRange)
	resp, err := tmc.cli.SendRequest(&node.GetChunk{
		Id: fileRange.id,
	})
	if err != nil {
		// try another server
		if len(fileRange.servers) == 0 {
			return err
		}
		fileRange.servers = fileRange.servers[1:]
		tmc.mgr.AddTaskContext(fileRange.servers[0], fileRange)
		return nil
	}
	chunk := resp.(*node.Chunk)
	return tmc.file.WritePart(fileRange.start, fileRange.end, chunk.Data)
}

// DownloadFile downloads the file from the remote file system
func (ctx *Context) DownloadFile(key string, filename string) error {
	pfile, err := pfile.NewPFile(filename, pfile.WritePart)
	if err != nil {
		return err
	}
	defer pfile.Close()

	resp, err := ctx.cli.SendRequest(&client.GetRequest{
		Key: key,
	})
	if err != nil {
		return err
	}

	respMsg := resp.(*client.GetResponse)
	allServers := gatherAllServers(respMsg.Chunks)

	taskMgr := taskmgr.NewTaskManager(ctx.handleDownload)
	defer closeClientsFromTaskMgr(taskMgr, allServers)
	for _, server := range allServers {
		conn, err := net.Dial("tcp", server)
		if err != nil {
			return err
		}
		taskMgr.AddContext(server, &taskMgrContext{
			file: pfile,
			cli:  rpcClient.NewClient(conn),
			mgr:  taskMgr,
		})
	}

	rangeStart := uint64(0)
	for _, chunk := range respMsg.Chunks {
		end := rangeStart + chunk.Size
		fileRange := fileRange{
			start:   int64(rangeStart),
			end:     int64(end),
			id:      chunk.Id,
			servers: chunk.Replicas,
		}
		rangeStart = end
		taskMgr.AddTaskContext(chunk.Replicas[0], &fileRange)
	}

	err = taskMgr.Run()
	if err == nil {
		fmt.Printf("Download file '%s' to '%s' successfully\n", key, filename)
	}
	return err
}

// DeleteFile deletes the file from the remote file system
func (ctx *Context) DeleteFile(key string) error {
	resp, err := ctx.cli.SendRequest(&client.DeleteRequest{
		Key: key,
	})
	if err != nil {
		return err
	}
	respMsg := resp.(*client.DeleteResponse)
	if respMsg.IsCandiate {
		fmt.Printf("tmp file '%s' deleted\n", key)
	} else {
		fmt.Printf("file '%s' deleted\n", key)
	}
	return nil
}

// ListFiles lists the files from the remote file system
func (ctx *Context) ListFiles(prefix string) error {
	resp, err := ctx.cli.SendRequest(&client.ListRequest{
		Prefix: prefix,
	})
	if err != nil {
		return err
	}
	respMsg := resp.(*client.ListResponse)
	fmt.Printf("Total files: %d\n", len(respMsg.Files))
	for _, file := range respMsg.Files {
		fmt.Printf("%d\t%s\t%s\n", file.Size, file.CreatedAt.AsTime().Local().Format(time.DateTime), file.Key)
	}
	return nil
}

// ListNodes lists the nodes from the remote file system
func (ctx Context) ListNodes(nodeId *string) error {
	resp, err := ctx.cli.SendRequest(&client.NodeInfoRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return err
	}
	respMsg := resp.(*client.NodeInfoResponse)
	if nodeId == nil {
		fmt.Printf("Total nodes: %d\n", len(respMsg.Nodes))
	}
	fmt.Printf("%15s%15s%15s\n", "node id", "free space", "total space")
	for _, node := range respMsg.Nodes {
		fmt.Printf("%15s%15d%15d\n", node.Id, node.Freespace, node.Totalspace)
	}
	return nil
}
