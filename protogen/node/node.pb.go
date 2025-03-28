// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: protos/node.proto

package node

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metadata struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Totalspace    uint64                 `protobuf:"varint,2,opt,name=totalspace,proto3" json:"totalspace,omitempty"`
	Freespace     uint64                 `protobuf:"varint,3,opt,name=freespace,proto3" json:"freespace,omitempty"`
	Chunks        []uint64               `protobuf:"varint,4,rep,packed,name=chunks,proto3" json:"chunks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	mi := &file_protos_node_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metadata) GetTotalspace() uint64 {
	if x != nil {
		return x.Totalspace
	}
	return 0
}

func (x *Metadata) GetFreespace() uint64 {
	if x != nil {
		return x.Freespace
	}
	return 0
}

func (x *Metadata) GetChunks() []uint64 {
	if x != nil {
		return x.Chunks
	}
	return nil
}

type Chunk struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Data  []byte                 `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// hash will be provided by the storage node
	Hash          *string `protobuf:"bytes,3,opt,name=hash,proto3,oneof" json:"hash,omitempty"`
	Size          *uint64 `protobuf:"varint,4,opt,name=size,proto3,oneof" json:"size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Chunk) Reset() {
	*x = Chunk{}
	mi := &file_protos_node_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chunk) ProtoMessage() {}

func (x *Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chunk.ProtoReflect.Descriptor instead.
func (*Chunk) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{1}
}

func (x *Chunk) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Chunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Chunk) GetHash() string {
	if x != nil && x.Hash != nil {
		return *x.Hash
	}
	return ""
}

func (x *Chunk) GetSize() uint64 {
	if x != nil && x.Size != nil {
		return *x.Size
	}
	return 0
}

type ChunkUploaded struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Hash          string                 `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Size          uint64                 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Replicas      []string               `protobuf:"bytes,4,rep,name=replicas,proto3" json:"replicas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChunkUploaded) Reset() {
	*x = ChunkUploaded{}
	mi := &file_protos_node_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChunkUploaded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChunkUploaded) ProtoMessage() {}

func (x *ChunkUploaded) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChunkUploaded.ProtoReflect.Descriptor instead.
func (*ChunkUploaded) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{2}
}

func (x *ChunkUploaded) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ChunkUploaded) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *ChunkUploaded) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ChunkUploaded) GetReplicas() []string {
	if x != nil {
		return x.Replicas
	}
	return nil
}

type Register struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Metadata      *Metadata              `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	ListenAddr    string                 `protobuf:"bytes,2,opt,name=listenAddr,proto3" json:"listenAddr,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Register) Reset() {
	*x = Register{}
	mi := &file_protos_node_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Register) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Register) ProtoMessage() {}

func (x *Register) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Register.ProtoReflect.Descriptor instead.
func (*Register) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{3}
}

func (x *Register) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Register) GetListenAddr() string {
	if x != nil {
		return x.ListenAddr
	}
	return ""
}

type Heartbeat struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Metadata      *Metadata              `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Heartbeat) Reset() {
	*x = Heartbeat{}
	mi := &file_protos_node_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Heartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat) ProtoMessage() {}

func (x *Heartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat.ProtoReflect.Descriptor instead.
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{4}
}

func (x *Heartbeat) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type HeartbeatResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RemoveChunks  []uint64               `protobuf:"varint,1,rep,packed,name=removeChunks,proto3" json:"removeChunks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HeartbeatResponse) Reset() {
	*x = HeartbeatResponse{}
	mi := &file_protos_node_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeartbeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatResponse) ProtoMessage() {}

func (x *HeartbeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatResponse.ProtoReflect.Descriptor instead.
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{5}
}

func (x *HeartbeatResponse) GetRemoveChunks() []uint64 {
	if x != nil {
		return x.RemoveChunks
	}
	return nil
}

type GetChunk struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetChunk) Reset() {
	*x = GetChunk{}
	mi := &file_protos_node_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChunk) ProtoMessage() {}

func (x *GetChunk) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChunk.ProtoReflect.Descriptor instead.
func (*GetChunk) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{6}
}

func (x *GetChunk) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type PutChunk struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Chunk *Chunk                 `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
	// replicas contains the list of node addresses to replicate the chunk to
	Replicas      []string `protobuf:"bytes,2,rep,name=replicas,proto3" json:"replicas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PutChunk) Reset() {
	*x = PutChunk{}
	mi := &file_protos_node_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PutChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutChunk) ProtoMessage() {}

func (x *PutChunk) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutChunk.ProtoReflect.Descriptor instead.
func (*PutChunk) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{7}
}

func (x *PutChunk) GetChunk() *Chunk {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *PutChunk) GetReplicas() []string {
	if x != nil {
		return x.Replicas
	}
	return nil
}

type DeleteChunk struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteChunk) Reset() {
	*x = DeleteChunk{}
	mi := &file_protos_node_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteChunk) ProtoMessage() {}

func (x *DeleteChunk) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteChunk.ProtoReflect.Descriptor instead.
func (*DeleteChunk) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteChunk) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type SendChunk struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Replicas      []string               `protobuf:"bytes,2,rep,name=replicas,proto3" json:"replicas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendChunk) Reset() {
	*x = SendChunk{}
	mi := &file_protos_node_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendChunk) ProtoMessage() {}

func (x *SendChunk) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendChunk.ProtoReflect.Descriptor instead.
func (*SendChunk) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{9}
}

func (x *SendChunk) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SendChunk) GetReplicas() []string {
	if x != nil {
		return x.Replicas
	}
	return nil
}

type SendChunks struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Chunks        []*SendChunk           `protobuf:"bytes,1,rep,name=chunks,proto3" json:"chunks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendChunks) Reset() {
	*x = SendChunks{}
	mi := &file_protos_node_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendChunks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendChunks) ProtoMessage() {}

func (x *SendChunks) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendChunks.ProtoReflect.Descriptor instead.
func (*SendChunks) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{10}
}

func (x *SendChunks) GetChunks() []*SendChunk {
	if x != nil {
		return x.Chunks
	}
	return nil
}

type SendChunksResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Chunks        []*SendChunkResponse   `protobuf:"bytes,1,rep,name=chunks,proto3" json:"chunks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendChunksResponse) Reset() {
	*x = SendChunksResponse{}
	mi := &file_protos_node_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendChunksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendChunksResponse) ProtoMessage() {}

func (x *SendChunksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendChunksResponse.ProtoReflect.Descriptor instead.
func (*SendChunksResponse) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{11}
}

func (x *SendChunksResponse) GetChunks() []*SendChunkResponse {
	if x != nil {
		return x.Chunks
	}
	return nil
}

type SendChunkResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Replicas      []string               `protobuf:"bytes,1,rep,name=replicas,proto3" json:"replicas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendChunkResponse) Reset() {
	*x = SendChunkResponse{}
	mi := &file_protos_node_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendChunkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendChunkResponse) ProtoMessage() {}

func (x *SendChunkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendChunkResponse.ProtoReflect.Descriptor instead.
func (*SendChunkResponse) Descriptor() ([]byte, []int) {
	return file_protos_node_proto_rawDescGZIP(), []int{12}
}

func (x *SendChunkResponse) GetReplicas() []string {
	if x != nil {
		return x.Replicas
	}
	return nil
}

var File_protos_node_proto protoreflect.FileDescriptor

var file_protos_node_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x66, 0x72, 0x65, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x66, 0x72, 0x65, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x04, 0x52, 0x06, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x73, 0x22, 0x6f, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x17, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x48, 0x01, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x63, 0x0a, 0x0d, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x22, 0x51, 0x0a, 0x08, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1e,
	0x0a, 0x0a, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x22, 0x32,
	0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x25, 0x0a, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x37, 0x0a, 0x11, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0c, 0x72,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x22, 0x1a, 0x0a, 0x08, 0x47,
	0x65, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x44, 0x0a, 0x08, 0x50, 0x75, 0x74, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x12, 0x1c, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x06, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x22, 0x1d, 0x0a,
	0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x37, 0x0a, 0x09,
	0x53, 0x65, 0x6e, 0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x73, 0x22, 0x30, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x68, 0x75,
	0x6e, 0x6b, 0x73, 0x12, 0x22, 0x0a, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52,
	0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x22, 0x40, 0x0a, 0x12, 0x53, 0x65, 0x6e, 0x64, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a,
	0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x52, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x22, 0x2f, 0x0a, 0x11, 0x53, 0x65, 0x6e,
	0x64, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x42, 0x13, 0x5a, 0x11, 0x64, 0x66,
	0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_protos_node_proto_rawDescOnce sync.Once
	file_protos_node_proto_rawDescData []byte
)

func file_protos_node_proto_rawDescGZIP() []byte {
	file_protos_node_proto_rawDescOnce.Do(func() {
		file_protos_node_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protos_node_proto_rawDesc), len(file_protos_node_proto_rawDesc)))
	})
	return file_protos_node_proto_rawDescData
}

var file_protos_node_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_protos_node_proto_goTypes = []any{
	(*Metadata)(nil),           // 0: Metadata
	(*Chunk)(nil),              // 1: Chunk
	(*ChunkUploaded)(nil),      // 2: ChunkUploaded
	(*Register)(nil),           // 3: Register
	(*Heartbeat)(nil),          // 4: Heartbeat
	(*HeartbeatResponse)(nil),  // 5: HeartbeatResponse
	(*GetChunk)(nil),           // 6: GetChunk
	(*PutChunk)(nil),           // 7: PutChunk
	(*DeleteChunk)(nil),        // 8: DeleteChunk
	(*SendChunk)(nil),          // 9: SendChunk
	(*SendChunks)(nil),         // 10: SendChunks
	(*SendChunksResponse)(nil), // 11: SendChunksResponse
	(*SendChunkResponse)(nil),  // 12: SendChunkResponse
}
var file_protos_node_proto_depIdxs = []int32{
	0,  // 0: Register.metadata:type_name -> Metadata
	0,  // 1: Heartbeat.metadata:type_name -> Metadata
	1,  // 2: PutChunk.chunk:type_name -> Chunk
	9,  // 3: SendChunks.chunks:type_name -> SendChunk
	12, // 4: SendChunksResponse.chunks:type_name -> SendChunkResponse
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_protos_node_proto_init() }
func file_protos_node_proto_init() {
	if File_protos_node_proto != nil {
		return
	}
	file_protos_node_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protos_node_proto_rawDesc), len(file_protos_node_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_node_proto_goTypes,
		DependencyIndexes: file_protos_node_proto_depIdxs,
		MessageInfos:      file_protos_node_proto_msgTypes,
	}.Build()
	File_protos_node_proto = out.File
	file_protos_node_proto_goTypes = nil
	file_protos_node_proto_depIdxs = nil
}
