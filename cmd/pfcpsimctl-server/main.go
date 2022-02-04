/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/internal/server"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// Values for mock-up4 environment
	defaultGNodeBAddress = "198.18.0.10"
	defaultUeAddressPool = "17.0.0.0/24"

	defaultUpfN3Address = "198.18.0.1"

	listenAddress = "0.0.0.0:54321" //TODO make address configurable
)

var (
	remotePeerAddress *string
	upfAddress        *string
	nodeBAddress      *string

	ueAddressPool *string
)

func getLocalAddress() (net.IP, error) {
	// cmd to run for darwin platforms
	cmd := "route -n get default | grep 'interface:' | grep -o '[^ ]*$'"

	if runtime.GOOS != "darwin" {
		// assuming linux platform
		cmd = "route | grep '^default' | grep -o '[^ ]*$'"
	}

	cmdOutput, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	interfaceName := strings.TrimSuffix(string(cmdOutput[:]), "\n")

	itf, _ := net.InterfaceByName(interfaceName)
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if v.IP.To4() != nil { //Verify if IP is IPV4
				ip = v.IP
			}
		}
	}

	if ip != nil {
		return ip, nil
	}

	return nil, fmt.Errorf("could not find interface: %v", interfaceName)
}

func startApiServer(apiDoneChannel chan bool, group *sync.WaitGroup) {
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("APIServer failed to listen: %v", err)
	}

	lAddr, err := getLocalAddress()
	if err != nil {
		log.Fatalf("Could not retrieve local address: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Initialize server
	pfcpServer, err := server.NewPFCPSimServer(lAddr, net.ParseIP(*remotePeerAddress))
	if err != nil {
		log.Fatalf("Could not create PFCPSimServer: %v", err)
	}

	pb.RegisterPFCPSimServer(grpcServer, pfcpServer)

	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to listed: %v", err)
		}
	}()

	log.Infof("Server listening on %v", listenAddress)

	x := <-apiDoneChannel
	if x {
		// if the API channel is closed, stop the gRPC server
		grpcServer.Stop()
		log.Warnf("Stopping API gRPC server")
	}

	group.Done()
}

func main() {
	// TODO make a config file instead of flag (?)

	remotePeerAddress = getopt.StringLong("remote-peer-address", 'r', "127.0.0.1", "Address or hostname of the remote peer (PFCP Agent)")
	upfAddress = getopt.StringLong("upf-address", 'u', defaultUpfN3Address, "Address of the UPF")
	ueAddressPool = getopt.StringLong("ue-address-pool", 'e', defaultUeAddressPool, "The IPv4 CIDR prefix from which UE addresses will be generated, incrementally")
	NodeBAddr := getopt.StringLong("nodeb-address", 'g', defaultGNodeBAddress, "The IPv4 of (g/e)NodeBAddress")

	optHelp := getopt.BoolLong("help", 0, "Help")

	getopt.Parse()
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	// Flag checks and validations
	if net.ParseIP(*NodeBAddr) == nil {
		log.Fatalf("Could not retrieve IP address of (g/e)NodeB")
	}

	if net.ParseIP(*remotePeerAddress) == nil {
		_, err := net.LookupHost(*remotePeerAddress)
		if err != nil {
			log.Fatalf("Could not retrieve hostname or address for remote peer: %s", *remotePeerAddress)
		}
	}

	if net.ParseIP(*upfAddress) == nil {
		log.Fatalf("Error while parsing UPF address")
	}

	_, _, err := net.ParseCIDR(*ueAddressPool)
	if err != nil {
		log.Fatalf("Could not parse ue address pool: %v", err)
	}

	// control channels, they are only closed when the goroutine needs to be terminated
	apiDoneChannel := make(chan bool)

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
