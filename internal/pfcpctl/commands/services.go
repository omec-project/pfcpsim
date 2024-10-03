// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	"context"

	"github.com/jessevdk/go-flags"
	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/logger"
)

type (
	associate                struct{}
	disassociate             struct{}
	configureRemoteAddresses struct {
		RemotePeerAddress  string `short:"r" long:"remote-peer-addr" default:"" description:"The remote PFCP agent address."`
		N3InterfaceAddress string `short:"n" long:"n3-addr" default:"" description:"UPF's N3 IP address"`
	}
)

type serviceOptions struct {
	Associate    associate                `command:"associate"`
	Disassociate disassociate             `command:"disassociate"`
	Configure    configureRemoteAddresses `command:"configure"`
}

func RegisterServiceCommands(parser *flags.Parser) {
	_, err := parser.AddCommand("service", "configure pfcpsim", "Command to configure pfcpsim", &serviceOptions{})
	if err != nil {
		logger.PfcpsimLog.Warnln(err)
	}
}

func (c *configureRemoteAddresses) Execute(args []string) error {
	client := connect()

	defer disconnect()

	res, err := client.Configure(context.Background(), &pb.ConfigureRequest{
		UpfN3Address:      c.N3InterfaceAddress,
		RemotePeerAddress: c.RemotePeerAddress,
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while configuring remote addresses: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)

	return nil
}

func (c *associate) Execute(args []string) error {
	client := connect()

	defer disconnect()

	res, err := client.Associate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while associating: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)

	return nil
}

func (c *disassociate) Execute(args []string) error {
	client := connect()

	defer disconnect()

	res, err := client.Disassociate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while disassociating: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)

	return nil
}
