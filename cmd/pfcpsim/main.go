/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/internal/pfcpsim"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	// Values for mock-up4 environment
	defaultGNodeBAddress = "198.18.0.10"
	defaultUeAddressPool = "17.0.0.0/24"

	defaultUpfN3Address = "198.18.0.1"

	defaultgRPCServerPort = "54321"
)

//var (
//	remotePeerAddress *string
//	upfAddress        *string
//	nodeBAddress      *string
//
//	ueAddressPool *string
//)

func startServer(apiDoneChannel chan bool, port string, group *sync.WaitGroup) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatalf("API gRPC Server failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterPFCPSimServer(grpcServer, &pfcpsim.ConcretePFCPSimServer{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to listed: %v", err)
		}
	}()

	log.Infof("Server listening on port %v", port)

	x := <-apiDoneChannel
	if x {
		// if the API channel is closed, stop the gRPC pfcpsim
		grpcServer.Stop()
		log.Warnf("Stopping API gRPC pfcpsim")
	}

	group.Done()
}

func main() {
	//remotePeerAddress = getopt.StringLong("remote-peer-address", 'r', "127.0.0.1", "Address or hostname of the remote peer (PFCP Agent)")
	//upfAddress = getopt.StringLong("upf-address", 'u', defaultUpfN3Address, "Address of the UPF")
	//ueAddressPool = getopt.StringLong("ue-address-pool", 'e', defaultUeAddressPool, "The IPv4 CIDR prefix from which UE addresses will be generated, incrementally")
	//nodeBAddress = getopt.StringLong("nodeb-address", 'g', defaultGNodeBAddress, "The IPv4 of (g/e)NodeBAddress")

	port := getopt.StringLong("port", 'p', defaultgRPCServerPort, "the gRPC Server port to listen")

	optHelp := getopt.BoolLong("help", 0, "Help")

	getopt.Parse()
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	// control channels, they are only closed when the goroutine needs to be terminated
	doneChannel := make(chan bool)

	sigs := make(chan os.Signal, 1)
	// stop API servers on SIGTERM
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		close(doneChannel)
	}()

	wg := sync.WaitGroup{}
	wg.Add(4)

	go startServer(doneChannel, *port, &wg)
	log.Debugf("Started API gRPC Service")

	wg.Wait()

	defer func() {
		log.Info("PFCP Simulator shutting down")
	}()
}
