/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package client

import (
	"context"

	pfcpsimctl "github.com/omec-project/pfcpsim/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PFCPSimClient struct {
}

func Connect() (pfcpsimctl.PFCPSimClient, *grpc.ClientConn) {
	serverAddress := ":9950" //TODO make this configurable

	// Create an insecure gRPC Channel
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect to %v: %v", serverAddress, err)
		return nil, conn
	}

	return pfcpsimctl.NewPFCPSimClient(conn), conn
}

func (P PFCPSimClient) SetLogLevel(ctx context.Context, in *pfcpsimctl.LogLevel, opts ...grpc.CallOption) (*pfcpsimctl.LogLevel, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimClient) StopgRPCServer(ctx context.Context, in *pfcpsimctl.Empty, opts ...grpc.CallOption) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimClient) StartgRPCServer(ctx context.Context, in *pfcpsimctl.Empty, opts ...grpc.CallOption) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimClient) Associate(ctx context.Context, in *pfcpsimctl.AssociateRequest, opts ...grpc.CallOption) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimClient) Disassociate(ctx context.Context, in *pfcpsimctl.DisassociateRequest, opts ...grpc.CallOption) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimClient) CreateSession(ctx context.Context, in *pfcpsimctl.CreateSessionRequest, opts ...grpc.CallOption) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimClient) ModifySession(ctx context.Context, in *pfcpsimctl.ModifySessionRequest, opts ...grpc.CallOption) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimClient) DeleteSession(ctx context.Context, in *pfcpsimctl.DeleteSessionRequest, opts ...grpc.CallOption) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}
