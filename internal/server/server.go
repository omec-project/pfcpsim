/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"context"
	"net"
	"net/http"

	pfcpsimctl "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	log "github.com/sirupsen/logrus"
	ieLib "github.com/wmnsk/go-pfcp/ie"
)

type pfcpClientContext struct {
	session *pfcpsim.PFCPSession

	pdrs []*ieLib.IE
	fars []*ieLib.IE
	qers []*ieLib.IE
}

// PFCPSimServer implements the Protobuf methods and keeps a connection to a remote PFCP Agent peer.
type PFCPSimServer struct {
	// Emulates 5G SMF/ 4G SGW
	client *pfcpsim.PFCPClient

	activeSessions []*pfcpClientContext
}

func NewPFCPSimServer(localAddress net.IP, remotePeerAddr net.IP) (*PFCPSimServer, error) {
	cl := pfcpsim.NewPFCPClient(localAddress.String())
	err := cl.ConnectN4(remotePeerAddr.String())

	log.Infof("pfcp client successfully connected to: %v", localAddress.String()) // DEBUG remove

	if err != nil {
		return nil, err
	}

	return &PFCPSimServer{
		client:         cl,
		activeSessions: make([]*pfcpClientContext, 0),
	}, nil
}

func (P PFCPSimServer) SetLogLevel(ctx context.Context, level *pfcpsimctl.LogLevel) (*pfcpsimctl.LogLevel, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) Associate(ctx context.Context, empty *pfcpsimctl.EmptyRequest) (*pfcpsimctl.Response, error) {
	err := P.client.SetupAssociation()
	if err != nil {
		return nil, err
	}

	return &pfcpsimctl.Response{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}, nil
}

func (P PFCPSimServer) Disassociate(ctx context.Context, empty *pfcpsimctl.EmptyRequest) (*pfcpsimctl.Response, error) {
	err := P.client.TeardownAssociation()
	if err != nil {
		return nil, err
	}

	return &pfcpsimctl.Response{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}, nil
}

func (P PFCPSimServer) CreateSession(ctx context.Context, empty *pfcpsimctl.EmptyRequest) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) ModifySession(ctx context.Context, empty *pfcpsimctl.EmptyRequest) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P PFCPSimServer) DeleteSession(ctx context.Context, empty *pfcpsimctl.EmptyRequest) (*pfcpsimctl.Response, error) {
	//TODO implement me
	panic("implement me")
}
