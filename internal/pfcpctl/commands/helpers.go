// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/internal/pfcpctl/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

func connect() pb.PFCPSimClient {
	// Create an insecure gRPC Channel
	conn, err := grpc.Dial(config.GlobalConfig.Server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error dialing %v: %v", config.GlobalConfig.Server, err)
	}

	return pb.NewPFCPSimClient(conn)
}

func disconnect() {
	if conn != nil {
		conn.Close()
	}
}
