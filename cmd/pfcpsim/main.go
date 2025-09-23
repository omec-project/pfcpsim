// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/internal/pfcpsim"
	"github.com/omec-project/pfcpsim/logger"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
)

const (
	defaultgRPCServerPort = "54321"
)

func startServer(apiDoneChannel chan bool, iFace string, port string, group *sync.WaitGroup) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		logger.PfcpsimLog.Fatalf("api gRPC Server failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterPFCPSimServer(grpcServer, pfcpsim.NewPFCPSimService(iFace))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.PfcpsimLog.Fatalf("failed to listed: %v", err)
		}
	}()

	logger.PfcpsimLog.Infoln("server listening on port", port)

	x := <-apiDoneChannel
	if x {
		// if the API channel is closed, stop the gRPC pfcpsim
		grpcServer.Stop()
		logger.PfcpsimLog.Warnln("stopping API gRPC pfcpsim")
	}

	group.Done()
}

func main() {
	app := &cli.Command{}
	app.Name = "pfcpsim"
	app.Usage = "./pfcpsim --interface <interface_name> --port <gRPC_server_port>"
	app.Flags = getCliFlags()
	app.Action = action

	logger.PfcpsimLog.Infof("app name: %s", app.Name)

	if err := app.Run(context.Background(), os.Args); err != nil {
		logger.PfcpsimLog.Fatalf("application error: %+v", err)
	}
}

func action(ctx context.Context, c *cli.Command) error {
	port := c.String("port")
	iFaceName := c.String("interface")

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
	wg.Add(1)

	go startServer(doneChannel, iFaceName, port, &wg)
	logger.PfcpsimLog.Debugln("started API gRPC Service")

	wg.Wait()

	logger.PfcpsimLog.Infoln("pfcp Simulator shutting down")
	return nil
}

func getCliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   defaultgRPCServerPort,
			Usage:   "the gRPC Server port to listen",
		},
		&cli.StringFlag{
			Name:    "interface",
			Aliases: []string{"i"},
			Value:   "",
			Usage:   "Defines the local address. If left blank, the IP will be taken from the first non-loopback interface",
		},
	}
}
