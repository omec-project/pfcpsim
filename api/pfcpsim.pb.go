// SPDX-License-Identifier: Apache-2.0
//Copyright 2022-present Open Networking Foundation

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.14.0
// source: pfcpsim.proto

package api

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type CreateSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *CreateSessionRequest) Reset() {
	*x = CreateSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsim_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSessionRequest) ProtoMessage() {}

func (x *CreateSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsim_proto_msgTypes[0]
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
	return file_pfcpsim_proto_rawDescGZIP(), []int{0}
}

func (x *CreateSessionRequest) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ModifySessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ModifySessionRequest) Reset() {
	*x = ModifySessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsim_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModifySessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModifySessionRequest) ProtoMessage() {}

func (x *ModifySessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsim_proto_msgTypes[1]
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
	return file_pfcpsim_proto_rawDescGZIP(), []int{1}
}

func (x *ModifySessionRequest) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type DeleteSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *DeleteSessionRequest) Reset() {
	*x = DeleteSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsim_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSessionRequest) ProtoMessage() {}

func (x *DeleteSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsim_proto_msgTypes[2]
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
	return file_pfcpsim_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteSessionRequest) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type EmptyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyRequest) Reset() {
	*x = EmptyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pfcpsim_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyRequest) ProtoMessage() {}

func (x *EmptyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsim_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyRequest.ProtoReflect.Descriptor instead.
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return file_pfcpsim_proto_rawDescGZIP(), []int{3}
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
		mi := &file_pfcpsim_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_pfcpsim_proto_msgTypes[4]
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
	return file_pfcpsim_proto_rawDescGZIP(), []int{4}
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

var File_pfcpsim_proto protoreflect.FileDescriptor

var file_pfcpsim_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x66, 0x63, 0x70, 0x73, 0x69, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x61, 0x70, 0x69, 0x22, 0x2c, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x2c, 0x0a, 0x14, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0x2c, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x0e,
	0x0a, 0x0c, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x45,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x91, 0x03, 0x0a, 0x07, 0x50, 0x46, 0x43, 0x50, 0x53, 0x69,
	0x6d, 0x12, 0x2f, 0x0a, 0x09, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x12, 0x11,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x32, 0x0a, 0x0c, 0x44, 0x69, 0x73, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61,
	0x74, 0x65, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0d, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x6f, 0x64, 0x69, 0x66,
	0x79, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x3b, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2f, 0x0a,
	0x09, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x72, 0x75, 0x70, 0x74, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x39,
	0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x50, 0x65, 0x65, 0x72, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x3b, 0x61,
	0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pfcpsim_proto_rawDescOnce sync.Once
	file_pfcpsim_proto_rawDescData = file_pfcpsim_proto_rawDesc
)

func file_pfcpsim_proto_rawDescGZIP() []byte {
	file_pfcpsim_proto_rawDescOnce.Do(func() {
		file_pfcpsim_proto_rawDescData = protoimpl.X.CompressGZIP(file_pfcpsim_proto_rawDescData)
	})
	return file_pfcpsim_proto_rawDescData
}

var file_pfcpsim_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pfcpsim_proto_goTypes = []interface{}{
	(*CreateSessionRequest)(nil), // 0: api.CreateSessionRequest
	(*ModifySessionRequest)(nil), // 1: api.ModifySessionRequest
	(*DeleteSessionRequest)(nil), // 2: api.DeleteSessionRequest
	(*EmptyRequest)(nil),         // 3: api.EmptyRequest
	(*Response)(nil),             // 4: api.Response
}
var file_pfcpsim_proto_depIdxs = []int32{
	3, // 0: api.PFCPSim.Associate:input_type -> api.EmptyRequest
	3, // 1: api.PFCPSim.Disassociate:input_type -> api.EmptyRequest
	0, // 2: api.PFCPSim.CreateSession:input_type -> api.CreateSessionRequest
	1, // 3: api.PFCPSim.ModifySession:input_type -> api.ModifySessionRequest
	2, // 4: api.PFCPSim.DeleteSession:input_type -> api.DeleteSessionRequest
	3, // 5: api.PFCPSim.Interrupt:input_type -> api.EmptyRequest
	3, // 6: api.PFCPSim.ConnectToRemotePeer:input_type -> api.EmptyRequest
	4, // 7: api.PFCPSim.Associate:output_type -> api.Response
	4, // 8: api.PFCPSim.Disassociate:output_type -> api.Response
	4, // 9: api.PFCPSim.CreateSession:output_type -> api.Response
	4, // 10: api.PFCPSim.ModifySession:output_type -> api.Response
	4, // 11: api.PFCPSim.DeleteSession:output_type -> api.Response
	4, // 12: api.PFCPSim.Interrupt:output_type -> api.Response
	4, // 13: api.PFCPSim.ConnectToRemotePeer:output_type -> api.Response
	7, // [7:14] is the sub-list for method output_type
	0, // [0:7] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pfcpsim_proto_init() }
func file_pfcpsim_proto_init() {
	if File_pfcpsim_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pfcpsim_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_pfcpsim_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_pfcpsim_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_pfcpsim_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyRequest); i {
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
		file_pfcpsim_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_pfcpsim_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pfcpsim_proto_goTypes,
		DependencyIndexes: file_pfcpsim_proto_depIdxs,
		MessageInfos:      file_pfcpsim_proto_msgTypes,
	}.Build()
	File_pfcpsim_proto = out.File
	file_pfcpsim_proto_rawDesc = nil
	file_pfcpsim_proto_goTypes = nil
	file_pfcpsim_proto_depIdxs = nil
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
	Associate(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error)
	Disassociate(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error)
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*Response, error)
	ModifySession(ctx context.Context, in *ModifySessionRequest, opts ...grpc.CallOption) (*Response, error)
	DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*Response, error)
	// Interrupt emulates pfcpsim crash
	Interrupt(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error)
	// Establish connection to remote peer (e.g. remote PFCP agent)
	ConnectToRemotePeer(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error)
}

type pFCPSimClient struct {
	cc grpc.ClientConnInterface
}

func NewPFCPSimClient(cc grpc.ClientConnInterface) PFCPSimClient {
	return &pFCPSimClient{cc}
}

func (c *pFCPSimClient) Associate(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/Associate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) Disassociate(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/Disassociate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/CreateSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) ModifySession(ctx context.Context, in *ModifySessionRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/ModifySession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/DeleteSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) Interrupt(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/Interrupt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pFCPSimClient) ConnectToRemotePeer(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/ConnectToRemotePeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PFCPSimServer is the server API for PFCPSim service.
type PFCPSimServer interface {
	Associate(context.Context, *EmptyRequest) (*Response, error)
	Disassociate(context.Context, *EmptyRequest) (*Response, error)
	CreateSession(context.Context, *CreateSessionRequest) (*Response, error)
	ModifySession(context.Context, *ModifySessionRequest) (*Response, error)
	DeleteSession(context.Context, *DeleteSessionRequest) (*Response, error)
	// Interrupt emulates pfcpsim crash
	Interrupt(context.Context, *EmptyRequest) (*Response, error)
	// Establish connection to remote peer (e.g. remote PFCP agent)
	ConnectToRemotePeer(context.Context, *EmptyRequest) (*Response, error)
}

// UnimplementedPFCPSimServer can be embedded to have forward compatible implementations.
type UnimplementedPFCPSimServer struct {
}

func (*UnimplementedPFCPSimServer) Associate(context.Context, *EmptyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Associate not implemented")
}
func (*UnimplementedPFCPSimServer) Disassociate(context.Context, *EmptyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disassociate not implemented")
}
func (*UnimplementedPFCPSimServer) CreateSession(context.Context, *CreateSessionRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (*UnimplementedPFCPSimServer) ModifySession(context.Context, *ModifySessionRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifySession not implemented")
}
func (*UnimplementedPFCPSimServer) DeleteSession(context.Context, *DeleteSessionRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSession not implemented")
}
func (*UnimplementedPFCPSimServer) Interrupt(context.Context, *EmptyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Interrupt not implemented")
}
func (*UnimplementedPFCPSimServer) ConnectToRemotePeer(context.Context, *EmptyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectToRemotePeer not implemented")
}

func RegisterPFCPSimServer(s *grpc.Server, srv PFCPSimServer) {
	s.RegisterService(&_PFCPSim_serviceDesc, srv)
}

func _PFCPSim_Associate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
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
		return srv.(PFCPSimServer).Associate(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_Disassociate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
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
		return srv.(PFCPSimServer).Disassociate(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSessionRequest)
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
		return srv.(PFCPSimServer).CreateSession(ctx, req.(*CreateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_ModifySession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifySessionRequest)
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
		return srv.(PFCPSimServer).ModifySession(ctx, req.(*ModifySessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_DeleteSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSessionRequest)
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
		return srv.(PFCPSimServer).DeleteSession(ctx, req.(*DeleteSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_Interrupt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).Interrupt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/Interrupt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).Interrupt(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PFCPSim_ConnectToRemotePeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).ConnectToRemotePeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/ConnectToRemotePeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).ConnectToRemotePeer(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PFCPSim_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.PFCPSim",
	HandlerType: (*PFCPSimServer)(nil),
	Methods: []grpc.MethodDesc{
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
		{
			MethodName: "Interrupt",
			Handler:    _PFCPSim_Interrupt_Handler,
		},
		{
			MethodName: "ConnectToRemotePeer",
			Handler:    _PFCPSim_ConnectToRemotePeer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pfcpsim.proto",
}
