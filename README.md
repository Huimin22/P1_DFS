# DFS Project

## Running `storage-node`

`storage-node` is the service responsible for storing files. It can be started with the following command:

```bash
./bin/storage-node -listenAddr localhost:9001 -storePath /path/to/store
```

### Parameters:
- `-listenAddr`: Specifies the address and port the storage node will listen on, default is `localhost:9001`.
- `-storePath`: The directory where files will be stored. Replace `/path/to/store` with your desired directory.

---

## Running `controller`

`controller` is the service that manages the nodes in the DFS network. To view the help information for the controller, run the following command:

```bash
./bin/controller
```

---

## Running `client`

`client` is the tool that interacts with the DFS network. You can view the help information for the client with:

```bash
./bin/client -h
```

---

## Common Commands

### Display Node Information

To list the information about all nodes:

```bash
./bin/client node
```

### List Files

To view the list of files stored in the DFS:

```bash
./bin/client list
```

### Upload File

To upload a file to the DFS system:

```bash
./bin/client upload /path/to/local/file.bin file_key
```

### Download File

To download a file from the DFS by its key:

```bash
./bin/client download file_key /path/to/downloaded/file.bin
```

### Delete File

To delete a file from the DFS by its key:

```bash
./bin/client delete file_key
```

### Verify File's MD5

To compute and verify the MD5 of a file:

```bash
md5 /path/to/file.bin
```

---

## Example Workflow

Assuming you have started the `storage-node` and `controller`, here is a typical usage workflow:

1. **Start Storage Node**:

   ```bash
   ./bin/storage-node -listenAddr localhost:9001 -storePath /path/to/store
   ```

2. **Check Node Information**:

   ```bash
   ./bin/client node
   ```

3. **Upload File**:

   ```bash
   ./bin/client upload /path/to/local/file.bin file_key
   ```

4. **Download File**:

   ```bash
   ./bin/client download file_key /path/to/downloaded/file.bin
   ```

5. **Delete File**:

   ```bash
   ./bin/client delete file_key
   ```

6. **MD5 Verification**:

   ```bash
   md5 /path/to/file.bin
   ```
