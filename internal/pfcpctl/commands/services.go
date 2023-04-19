// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	"context"

	"github.com/jessevdk/go-flags"
	pb "github.com/omec-project/pfcpsim/api"
	log "github.com/sirupsen/logrus"
)

type associate struct{}
type disassociate struct{}
type configureRemoteAddresses struct {
	RemotePeerAddress  string `short:"r" long:"remote-peer-addr" default:"" description:"The remote PFCP agent address."`
	N3InterfaceAddress string `short:"n" long:"n3-addr" default:"" description:"The IPv4 address of the UPF's N3 interface"`
}

type serviceOptions struct {
	Associate    associate                `command:"associate"`
	Disassociate disassociate             `command:"disassociate"`
	Configure    configureRemoteAddresses `command:"configure"`
}

func RegisterServiceCommands(parser *flags.Parser) {
	_, _ = parser.AddCommand("service", "configure pfcpsim", "Command to configure pfcpsim", &serviceOptions{})
}

func (c *configureRemoteAddresses) Execute(args []string) error {
	client := connect()

	defer disconnect()

	res, err := client.Configure(context.Background(), &pb.ConfigureRequest{
		UpfN3Address:      c.N3InterfaceAddress,
		RemotePeerAddress: c.RemotePeerAddress,
	})

	if err != nil {
		log.Fatalf("Error while configuring remote addresses: %v", err)
	}

	log.Info(res.Message)

	return nil
}

func (c *associate) Execute(args []string) error {
	client := connect()

	defer disconnect()

	res, err := client.Associate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("Error while associating: %v", err)
	}

	log.Infof(res.Message)

	return nil
}

func (c *disassociate) Execute(args []string) error {
	client := connect()

	defer disconnect()

	res, err := client.Disassociate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("Error while disassociating: %v", err)
	}

	log.Infof(res.Message)

	return nil
}
