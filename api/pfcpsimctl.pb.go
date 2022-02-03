// SPDX-License-Identifier: Apache-2.0
//Copyright 2022-present Open Networking Foundation

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.11.4
// source: pfcpsimctl.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogLevel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Level  string `protobuf:"bytes,1,opt,name=level,proto3" json:"level,omitempty"`
	Caller bool   `protobuf:"varint,2,opt,name=caller,proto3" json:"caller,omitempty"`
}

func (x *LogLevel) Reset() {
	*x = LogLevel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogLevel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogLevel) ProtoMessage() {}

func (x *LogLevel) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogLevel.ProtoReflect.Descriptor instead.
func (*LogLevel) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{0}
}

func (x *LogLevel) GetLevel() string {
	if x != nil {
		return x.Level
	}
	return ""
}

func (x *LogLevel) GetCaller() bool {
	if x != nil {
		return x.Caller
	}
	return false
}

type AssociateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AssociateRequest) Reset() {
	*x = AssociateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssociateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssociateRequest) ProtoMessage() {}

func (x *AssociateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssociateRequest.ProtoReflect.Descriptor instead.
func (*AssociateRequest) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{1}
}

type DisassociateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DisassociateRequest) Reset() {
	*x = DisassociateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DisassociateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DisassociateRequest) ProtoMessage() {}

func (x *DisassociateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DisassociateRequest.ProtoReflect.Descriptor instead.
func (*DisassociateRequest) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{2}
}

type CreateSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateSessionRequest) Reset() {
	*x = CreateSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSessionRequest) ProtoMessage() {}

func (x *CreateSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSessionRequest.ProtoReflect.Descriptor instead.
func (*CreateSessionRequest) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{3}
}

type ModifySessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ModifySessionRequest) Reset() {
	*x = ModifySessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModifySessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModifySessionRequest) ProtoMessage() {}

func (x *ModifySessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModifySessionRequest.ProtoReflect.Descriptor instead.
func (*ModifySessionRequest) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{4}
}

type DeleteSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteSessionRequest) Reset() {
	*x = DeleteSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSessionRequest) ProtoMessage() {}

func (x *DeleteSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSessionRequest.ProtoReflect.Descriptor instead.
func (*DeleteSessionRequest) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{5}
}

type InterruptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InterruptRequest) Reset() {
	*x = InterruptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterruptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterruptRequest) ProtoMessage() {}

func (x *InterruptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterruptRequest.ProtoReflect.Descriptor instead.
func (*InterruptRequest) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{6}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{7}
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Message    string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsimctl_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsimctl_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_pfcpsimctl_proto_rawDescGZIP(), []int{8}
}

func (x *Response) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_pfcpsimctl_proto protoreflect.FileDescriptor

var file_pfcpsimctl_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x66, 0x63, 0x70, 0x73, 0x69, 0x6d, 0x63, 0x74, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x22, 0x38, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x6c,
	0x6c, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x63, 0x61, 0x6c, 0x6c, 0x65,
	0x72, 0x22, 0x12, 0x0a, 0x10, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x15, 0x0a, 0x13, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f,
	0x63, 0x69, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x16, 0x0a, 0x14,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x16, 0x0a, 0x14, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x16, 0x0a, 0x14,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x12, 0x0a, 0x10, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x72, 0x75, 0x70,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x45, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xf8, 0x02, 0x0a, 0x07, 0x50, 0x46, 0x43,
	0x50, 0x53, 0x69, 0x6d, 0x12, 0x2d, 0x0a, 0x0b, 0x53, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x12, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76,
	0x65, 0x6c, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65,
	0x6c, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x0e, 0x53, 0x74, 0x6f, 0x70, 0x67, 0x52, 0x50, 0x43, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x2e, 0x0a, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x67, 0x52, 0x50, 0x43, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x28, 0x0a, 0x09, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x12,
	0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x0c,
	0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x0d, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0a, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x0d, 0x4d, 0x6f, 0x64, 0x69, 0x66,
	0x79, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pfcpsimctl_proto_rawDescOnce sync.Once
	file_pfcpsimctl_proto_rawDescData = file_pfcpsimctl_proto_rawDesc
)

func file_pfcpsimctl_proto_rawDescGZIP() []byte {
	file_pfcpsimctl_proto_rawDescOnce.Do(func() {
		file_pfcpsimctl_proto_rawDescData = protoimpl.X.CompressGZIP(file_pfcpsimctl_proto_rawDescData)
	})
	return file_pfcpsimctl_proto_rawDescData
}

var file_pfcpsimctl_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_pfcpsimctl_proto_goTypes = []interface{}{
	(*LogLevel)(nil),             // 0: api.LogLevel
	(*AssociateRequest)(nil),     // 1: api.AssociateRequest
	(*DisassociateRequest)(nil),  // 2: api.DisassociateRequest
	(*CreateSessionRequest)(nil), // 3: api.CreateSessionRequest
	(*ModifySessionRequest)(nil), // 4: api.ModifySessionRequest
	(*DeleteSessionRequest)(nil), // 5: api.DeleteSessionRequest
	(*InterruptRequest)(nil),     // 6: api.InterruptRequest
	(*Empty)(nil),                // 7: api.Empty
	(*Response)(nil),             // 8: api.Response
}
var file_pfcpsimctl_proto_depIdxs = []int32{
	0, // 0: api.PFCPSim.SetLogLevel:input_type -> api.LogLevel
	7, // 1: api.PFCPSim.StopgRPCServer:input_type -> api.Empty
	7, // 2: api.PFCPSim.StartgRPCServer:input_type -> api.Empty
	7, // 3: api.PFCPSim.Associate:input_type -> api.Empty
	7, // 4: api.PFCPSim.Disassociate:input_type -> api.Empty
	7, // 5: api.PFCPSim.CreateSession:input_type -> api.Empty
	7, // 6: api.PFCPSim.ModifySession:input_type -> api.Empty
	7, // 7: api.PFCPSim.DeleteSession:input_type -> api.Empty
	0, // 8: api.PFCPSim.SetLogLevel:output_type -> api.LogLevel
	8, // 9: api.PFCPSim.StopgRPCServer:output_type -> api.Response
	8, // 10: api.PFCPSim.StartgRPCServer:output_type -> api.Response
	8, // 11: api.PFCPSim.Associate:output_type -> api.Response
	8, // 12: api.PFCPSim.Disassociate:output_type -> api.Response
	8, // 13: api.PFCPSim.CreateSession:output_type -> api.Response
	8, // 14: api.PFCPSim.ModifySession:output_type -> api.Response
	8, // 15: api.PFCPSim.DeleteSession:output_type -> api.Response
	8, // [8:16] is the sub-list for method output_type
	0, // [0:8] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pfcpsimctl_proto_init() }
func file_pfcpsimctl_proto_init() {
	if File_pfcpsimctl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pfcpsimctl_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogLevel); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssociateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DisassociateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSessionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ModifySessionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteSessionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InterruptRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pfcpsimctl_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pfcpsimctl_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pfcpsimctl_proto_goTypes,
		DependencyIndexes: file_pfcpsimctl_proto_depIdxs,
		MessageInfos:      file_pfcpsimctl_proto_msgTypes,
	}.Build()
	File_pfcpsimctl_proto = out.File
	file_pfcpsimctl_proto_rawDesc = nil
	file_pfcpsimctl_proto_goTypes = nil
	file_pfcpsimctl_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PFCPSimClient is the client API for PFCPSim service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PFCPSimClient interface {
	// Set PFCPSim log level
	SetLogLevel(ctx context.Context, in *LogLevel, opts ...grpc.CallOption) (*LogLevel, error)
	StopgRPCServer(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	StartgRPCServer(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	Associate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	Disassociate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	CreateSession(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	ModifySession(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	DeleteSession(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
}

type pFCPSimClient struct {
	cc grpc.ClientConnInterface
}

func NewPFCPSimClient(cc grpc.ClientConnInterface) PFCPSimClient {
	return &pFCPSimClient{cc}
}

func (c *pFCPSimClient) SetLogLevel(ctx context.Context, in *LogLevel, opts ...grpc.CallOption) (*LogLevel, error) {
	out := new(LogLevel)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/SetLogLevel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) StopgRPCServer(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/StopgRPCServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) StartgRPCServer(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/StartgRPCServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) Associate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/Associate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) Disassociate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/Disassociate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) CreateSession(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/CreateSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) ModifySession(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/ModifySession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) DeleteSession(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/DeleteSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PFCPSimServer is the server API for PFCPSim service.
type PFCPSimServer interface {
	// Set PFCPSim log level
	SetLogLevel(context.Context, *LogLevel) (*LogLevel, error)
	StopgRPCServer(context.Context, *Empty) (*Response, error)
	StartgRPCServer(context.Context, *Empty) (*Response, error)
	Associate(context.Context, *Empty) (*Response, error)
	Disassociate(context.Context, *Empty) (*Response, error)
	CreateSession(context.Context, *Empty) (*Response, error)
	ModifySession(context.Context, *Empty) (*Response, error)
	DeleteSession(context.Context, *Empty) (*Response, error)
}

// UnimplementedPFCPSimServer can be embedded to have forward compatible implementations.
type UnimplementedPFCPSimServer struct {
}

func (*UnimplementedPFCPSimServer) SetLogLevel(context.Context, *LogLevel) (*LogLevel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLogLevel not implemented")
}
func (*UnimplementedPFCPSimServer) StopgRPCServer(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopgRPCServer not implemented")
}
func (*UnimplementedPFCPSimServer) StartgRPCServer(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartgRPCServer not implemented")
}
func (*UnimplementedPFCPSimServer) Associate(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Associate not implemented")
}
func (*UnimplementedPFCPSimServer) Disassociate(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disassociate not implemented")
}
func (*UnimplementedPFCPSimServer) CreateSession(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (*UnimplementedPFCPSimServer) ModifySession(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifySession not implemented")
}
func (*UnimplementedPFCPSimServer) DeleteSession(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSession not implemented")
}

func RegisterPFCPSimServer(s *grpc.Server, srv PFCPSimServer) {
	s.RegisterService(&_PFCPSim_serviceDesc, srv)
}

func _PFCPSim_SetLogLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogLevel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).SetLogLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/SetLogLevel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).SetLogLevel(ctx, req.(*LogLevel))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_StopgRPCServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).StopgRPCServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/StopgRPCServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).StopgRPCServer(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_StartgRPCServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).StartgRPCServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/StartgRPCServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).StartgRPCServer(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_Associate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).Associate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/Associate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).Associate(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_Disassociate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).Disassociate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/Disassociate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).Disassociate(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/CreateSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).CreateSession(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_ModifySession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).ModifySession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/ModifySession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).ModifySession(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_DeleteSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).DeleteSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/DeleteSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).DeleteSession(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _PFCPSim_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.PFCPSim",
	HandlerType: (*PFCPSimServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetLogLevel",
			Handler:    _PFCPSim_SetLogLevel_Handler,
		},
		{
			MethodName: "StopgRPCServer",
			Handler:    _PFCPSim_StopgRPCServer_Handler,
		},
		{
			MethodName: "StartgRPCServer",
			Handler:    _PFCPSim_StartgRPCServer_Handler,
		},
		{
			MethodName: "Associate",
			Handler:    _PFCPSim_Associate_Handler,
		},
		{
			MethodName: "Disassociate",
			Handler:    _PFCPSim_Disassociate_Handler,
		},
		{
			MethodName: "CreateSession",
			Handler:    _PFCPSim_CreateSession_Handler,
		},
		{
			MethodName: "ModifySession",
			Handler:    _PFCPSim_ModifySession_Handler,
		},
		{
			MethodName: "DeleteSession",
			Handler:    _PFCPSim_DeleteSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pfcpsimctl.proto",
}
