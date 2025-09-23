// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	"context"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/logger"
	"github.com/urfave/cli/v3"
)

func GetServiceCommands() *cli.Command {
	return &cli.Command{
		Name:  "service",
		Usage: "configure pfcpsim",
		Commands: []*cli.Command{
			{
				Name:  "associate",
				Usage: "Associate with remote PFCP agent",
				Action: func(ctx context.Context, c *cli.Command) error {
					return associateAction(ctx, c)
				},
			},
			{
				Name:  "disassociate",
				Usage: "Disassociate from remote PFCP agent",
				Action: func(ctx context.Context, c *cli.Command) error {
					return disassociateAction(ctx, c)
				},
			},
			{
				Name:  "configure",
				Usage: "Configure remote addresses",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "remote-peer-addr",
						Aliases: []string{"r"},
						Usage:   "The remote PFCP agent address",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "n3-addr",
						Aliases: []string{"n"},
						Usage:   "UPF's N3 IP address",
						Value:   "",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return configureAction(ctx, c)
				},
			},
		},
	}
}

func configureAction(ctx context.Context, c *cli.Command) error {
	client := connect()
	defer disconnect()

	remotePeerAddr := c.String("remote-peer-addr")
	n3Addr := c.String("n3-addr")

	res, err := client.Configure(ctx, &pb.ConfigureRequest{
		UpfN3Address:      n3Addr,
		RemotePeerAddress: remotePeerAddr,
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while configuring remote addresses: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)
	return nil
}

func associateAction(ctx context.Context, _ *cli.Command) error {
	client := connect()
	defer disconnect()

	res, err := client.Associate(ctx, &pb.EmptyRequest{})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while associating: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)
	return nil
}

func disassociateAction(ctx context.Context, _ *cli.Command) error {
	client := connect()
	defer disconnect()

	res, err := client.Disassociate(ctx, &pb.EmptyRequest{})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while disassociating: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)
	return nil
}
