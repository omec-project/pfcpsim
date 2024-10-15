// Copyright 2022-present Open Networking Foundation
// Copyright 2024-present Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConnInterface
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PFCPSimClient is the client API for PFCPSim service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PFCPSimClient interface {
	Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*Response, error)
	// Associate connects PFCPClient to remote peer and starts an association
	Associate(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error)
	// Disassociate perform teardown of association and disconnects from remote peer.
	Disassociate(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error)
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*Response, error)
	ModifySession(ctx context.Context, in *ModifySessionRequest, opts ...grpc.CallOption) (*Response, error)
	DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*Response, error)
}

type pFCPSimClient struct {
	cc grpc.ClientConnInterface
}

func NewPFCPSimClient(cc grpc.ClientConnInterface) PFCPSimClient {
	return &pFCPSimClient{cc}
}

func (c *pFCPSimClient) Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.PFCPSim/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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

// PFCPSimServer is the server API for PFCPSim service.
type PFCPSimServer interface {
	Configure(context.Context, *ConfigureRequest) (*Response, error)
	// Associate connects PFCPClient to remote peer and starts an association
	Associate(context.Context, *EmptyRequest) (*Response, error)
	// Disassociate perform teardown of association and disconnects from remote peer.
	Disassociate(context.Context, *EmptyRequest) (*Response, error)
	CreateSession(context.Context, *CreateSessionRequest) (*Response, error)
	ModifySession(context.Context, *ModifySessionRequest) (*Response, error)
	DeleteSession(context.Context, *DeleteSessionRequest) (*Response, error)
}

// UnimplementedPFCPSimServer can be embedded to have forward compatible implementations.
type UnimplementedPFCPSimServer struct{}

func (*UnimplementedPFCPSimServer) Configure(context.Context, *ConfigureRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
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

func RegisterPFCPSimServer(s *grpc.Server, srv PFCPSimServer) {
	s.RegisterService(&_PFCPSim_serviceDesc, srv)
}

func _PFCPSim_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PFCPSimServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PFCPSim/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PFCPSimServer).Configure(ctx, req.(*ConfigureRequest))
	}
	return interceptor(ctx, in, info, handler)
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

var _PFCPSim_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.PFCPSim",
	HandlerType: (*PFCPSimServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Configure",
			Handler:    _PFCPSim_Configure_Handler,
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
	Metadata: "pfcpsim.proto",
}
