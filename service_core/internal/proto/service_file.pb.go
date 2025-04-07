// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.1
// source: service_file.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type GetFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileId        uint32                 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileRequest) Reset() {
	*x = GetFileRequest{}
	mi := &file_service_file_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileRequest) ProtoMessage() {}

func (x *GetFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileRequest.ProtoReflect.Descriptor instead.
func (*GetFileRequest) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{0}
}

func (x *GetFileRequest) GetFileId() uint32 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type FileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DirectoryId   uint32                 `protobuf:"varint,2,opt,name=directory_id,json=directoryId,proto3" json:"directory_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Status        string                 `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	Version       int32                  `protobuf:"varint,5,opt,name=version,proto3" json:"version,omitempty"`
	Directory     *DirectoryResponse     `protobuf:"bytes,6,opt,name=directory,proto3" json:"directory,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileResponse) Reset() {
	*x = FileResponse{}
	mi := &file_service_file_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileResponse) ProtoMessage() {}

func (x *FileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileResponse.ProtoReflect.Descriptor instead.
func (*FileResponse) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{1}
}

func (x *FileResponse) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FileResponse) GetDirectoryId() uint32 {
	if x != nil {
		return x.DirectoryId
	}
	return 0
}

func (x *FileResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FileResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *FileResponse) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *FileResponse) GetDirectory() *DirectoryResponse {
	if x != nil {
		return x.Directory
	}
	return nil
}

type DirectoryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ParentPathId  uint32                 `protobuf:"varint,2,opt,name=parent_path_id,json=parentPathId,proto3" json:"parent_path_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Status        string                 `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	WorkflowId    uint32                 `protobuf:"varint,5,opt,name=workflow_id,json=workflowId,proto3" json:"workflow_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DirectoryResponse) Reset() {
	*x = DirectoryResponse{}
	mi := &file_service_file_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DirectoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirectoryResponse) ProtoMessage() {}

func (x *DirectoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirectoryResponse.ProtoReflect.Descriptor instead.
func (*DirectoryResponse) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{2}
}

func (x *DirectoryResponse) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DirectoryResponse) GetParentPathId() uint32 {
	if x != nil {
		return x.ParentPathId
	}
	return 0
}

func (x *DirectoryResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DirectoryResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *DirectoryResponse) GetWorkflowId() uint32 {
	if x != nil {
		return x.WorkflowId
	}
	return 0
}

type UpdateFileStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileId        uint32                 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateFileStatusRequest) Reset() {
	*x = UpdateFileStatusRequest{}
	mi := &file_service_file_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateFileStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateFileStatusRequest) ProtoMessage() {}

func (x *UpdateFileStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateFileStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdateFileStatusRequest) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateFileStatusRequest) GetFileId() uint32 {
	if x != nil {
		return x.FileId
	}
	return 0
}

func (x *UpdateFileStatusRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type GetFilesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileIds       []uint32               `protobuf:"varint,1,rep,packed,name=file_ids,json=fileIds,proto3" json:"file_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFilesRequest) Reset() {
	*x = GetFilesRequest{}
	mi := &file_service_file_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFilesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFilesRequest) ProtoMessage() {}

func (x *GetFilesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFilesRequest.ProtoReflect.Descriptor instead.
func (*GetFilesRequest) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{4}
}

func (x *GetFilesRequest) GetFileIds() []uint32 {
	if x != nil {
		return x.FileIds
	}
	return nil
}

type GetFilesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileNames     map[uint32]string      `protobuf:"bytes,1,rep,name=file_names,json=fileNames,proto3" json:"file_names,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFilesResponse) Reset() {
	*x = GetFilesResponse{}
	mi := &file_service_file_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFilesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFilesResponse) ProtoMessage() {}

func (x *GetFilesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFilesResponse.ProtoReflect.Descriptor instead.
func (*GetFilesResponse) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{5}
}

func (x *GetFilesResponse) GetFileNames() map[uint32]string {
	if x != nil {
		return x.FileNames
	}
	return nil
}

type CheckWorkflowRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WorkflowId    uint32                 `protobuf:"varint,1,opt,name=workflow_id,json=workflowId,proto3" json:"workflow_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckWorkflowRequest) Reset() {
	*x = CheckWorkflowRequest{}
	mi := &file_service_file_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckWorkflowRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckWorkflowRequest) ProtoMessage() {}

func (x *CheckWorkflowRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckWorkflowRequest.ProtoReflect.Descriptor instead.
func (*CheckWorkflowRequest) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{6}
}

func (x *CheckWorkflowRequest) GetWorkflowId() uint32 {
	if x != nil {
		return x.WorkflowId
	}
	return 0
}

type CheckWorkflowResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Exists        bool                   `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckWorkflowResponse) Reset() {
	*x = CheckWorkflowResponse{}
	mi := &file_service_file_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckWorkflowResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckWorkflowResponse) ProtoMessage() {}

func (x *CheckWorkflowResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_file_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckWorkflowResponse.ProtoReflect.Descriptor instead.
func (*CheckWorkflowResponse) Descriptor() ([]byte, []int) {
	return file_service_file_proto_rawDescGZIP(), []int{7}
}

func (x *CheckWorkflowResponse) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

var File_service_file_proto protoreflect.FileDescriptor

var file_service_file_proto_rawDesc = string([]byte{
	0x0a, 0x12, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x29, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65,
	0x49, 0x64, 0x22, 0xbe, 0x01, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x09,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x79, 0x22, 0x96, 0x01, 0x0a, 0x11, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x74, 0x68, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x77,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x49, 0x64, 0x22, 0x4a, 0x0a, 0x17,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x66,
	0x69, 0x6c, 0x65, 0x49, 0x64, 0x73, 0x22, 0x96, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x0a, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x1a, 0x3c, 0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x37, 0x0a, 0x14, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x49, 0x64, 0x22, 0x2f, 0x0a, 0x15, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x32, 0x9a, 0x02, 0x0a, 0x0b, 0x46, 0x69,
	0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x49, 0x44, 0x12, 0x14, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e,
	0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x49, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1d, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3d, 0x0a,
	0x0c, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x15, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x0d,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x12, 0x1a, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x17, 0x5a, 0x15, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2d, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_service_file_proto_rawDescOnce sync.Once
	file_service_file_proto_rawDescData []byte
)

func file_service_file_proto_rawDescGZIP() []byte {
	file_service_file_proto_rawDescOnce.Do(func() {
		file_service_file_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_service_file_proto_rawDesc), len(file_service_file_proto_rawDesc)))
	})
	return file_service_file_proto_rawDescData
}

var file_service_file_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_service_file_proto_goTypes = []any{
	(*GetFileRequest)(nil),          // 0: file.GetFileRequest
	(*FileResponse)(nil),            // 1: file.FileResponse
	(*DirectoryResponse)(nil),       // 2: file.DirectoryResponse
	(*UpdateFileStatusRequest)(nil), // 3: file.UpdateFileStatusRequest
	(*GetFilesRequest)(nil),         // 4: file.GetFilesRequest
	(*GetFilesResponse)(nil),        // 5: file.GetFilesResponse
	(*CheckWorkflowRequest)(nil),    // 6: file.CheckWorkflowRequest
	(*CheckWorkflowResponse)(nil),   // 7: file.CheckWorkflowResponse
	nil,                             // 8: file.GetFilesResponse.FileNamesEntry
	(*emptypb.Empty)(nil),           // 9: google.protobuf.Empty
}
var file_service_file_proto_depIdxs = []int32{
	2, // 0: file.FileResponse.directory:type_name -> file.DirectoryResponse
	8, // 1: file.GetFilesResponse.file_names:type_name -> file.GetFilesResponse.FileNamesEntry
	0, // 2: file.FileService.GetFileByID:input_type -> file.GetFileRequest
	3, // 3: file.FileService.UpdateFileStatus:input_type -> file.UpdateFileStatusRequest
	4, // 4: file.FileService.GetFilesInfo:input_type -> file.GetFilesRequest
	6, // 5: file.FileService.CheckWorkflow:input_type -> file.CheckWorkflowRequest
	1, // 6: file.FileService.GetFileByID:output_type -> file.FileResponse
	9, // 7: file.FileService.UpdateFileStatus:output_type -> google.protobuf.Empty
	5, // 8: file.FileService.GetFilesInfo:output_type -> file.GetFilesResponse
	7, // 9: file.FileService.CheckWorkflow:output_type -> file.CheckWorkflowResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_service_file_proto_init() }
func file_service_file_proto_init() {
	if File_service_file_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_service_file_proto_rawDesc), len(file_service_file_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_file_proto_goTypes,
		DependencyIndexes: file_service_file_proto_depIdxs,
		MessageInfos:      file_service_file_proto_msgTypes,
	}.Build()
	File_service_file_proto = out.File
	file_service_file_proto_goTypes = nil
	file_service_file_proto_depIdxs = nil
}
