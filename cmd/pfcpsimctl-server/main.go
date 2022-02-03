/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package pfcpsimctl_server

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pfcpsimctl "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/internal/server"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startApiServer(apiDoneChannel chan bool, group *sync.WaitGroup) {
	lis, err := net.Listen("tcp", "0.0.0.0")
	if err != nil {
		log.Fatalf("APIServer failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pfcpsimctl.RegisterPFCPSimServer(grpcServer, server.PFCPSimServer{})

	reflection.Register(grpcServer)

	go func() { _ = grpcServer.Serve(lis) }()

	x := <-apiDoneChannel
	if x {
		// if the API channel is closed, stop the gRPC server
		grpcServer.Stop()
		log.Warnf("Stopping API gRPC server")
	}

	group.Done()
}

func main() {

	// control channels, they are only closed when the goroutine needs to be terminated
	apiDoneChannel := make(chan bool)

	simClient := &pfcpsim.PFCPClient{}

	log.Debugf("Created new pfcpsim client %v", simClient)

	sigs := make(chan os.Signal, 1)
	// stop API servers on SIGTERM
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		close(apiDoneChannel)
	}()

	wg := sync.WaitGroup{}
	wg.Add(4)

	go startApiServer(apiDoneChannel, &wg)
	log.Debugf("Started APIService")

	wg.Wait()

	defer func() {
		log.Info("PFCP Simulator shutting down")
	}()
}
