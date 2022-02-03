/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"context"

	pfcpsimctl "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	ieLib "github.com/wmnsk/go-pfcp/ie"
)

type pfcpClientContext struct {
	session *pfcpsim.PFCPSession

	pdrs []*ieLib.IE
	fars []*ieLib.IE
	qers []*ieLib.IE
}

var ()

// PFCPSimServer implements the Protobuf methods and keeps a connection to a remote PFCP Agent peer.
type PFCPSimServer struct {
	// Emulates 5G SMF/ 4G SGW
	Client *pfcpsim.PFCPClient

	activeSessions []*pfcpClientContext
}

func NewPFCPSimServer(client *pfcpsim.PFCPClient) *PFCPSimServer {
	return &PFCPSimServer{
		Client:         client,
		activeSessions: make([]*pfcpClientContext, 0),
	}
}

func (P PFCPSimServer) SetLogLevel(ctx context.Context, level *pfcpsimctl.LogLevel) (*pfcpsimctl.LogLevel, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) StopgRPCServer(ctx context.Context, empty *pfcpsimctl.Empty) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) StartgRPCServer(ctx context.Context, empty *pfcpsimctl.Empty) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) Associate(ctx context.Context, empty *pfcpsimctl.Empty) (*pfcpsimctl.Response, error) {
	P.Client.ConnectN4()
}

func (P PFCPSimServer) Disassociate(ctx context.Context, empty *pfcpsimctl.Empty) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) CreateSession(ctx context.Context, empty *pfcpsimctl.Empty) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) ModifySession(ctx context.Context, empty *pfcpsimctl.Empty) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) DeleteSession(ctx context.Context, empty *pfcpsimctl.Empty) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}
