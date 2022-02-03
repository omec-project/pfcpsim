/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"context"

	pfcpsimctl "github.com/omec-project/pfcpsim/api"
)

type PFCPSimServer struct {
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
	//TODO implement me
	panic("implement me")
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
