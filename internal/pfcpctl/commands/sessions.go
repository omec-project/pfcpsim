// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	"context"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/logger"
	"github.com/urfave/cli/v3"
)

// getCommonFlags returns the common flags used by session commands
func getCommonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"c"},
			Value:   1,
			Usage:   "The number of sessions to create",
		},
		&cli.IntFlag{
			Name:    "baseID",
			Aliases: []string{"i"},
			Value:   1,
			Usage:   "The base ID to use",
		},
		&cli.StringFlag{
			Name:    "ue-pool",
			Aliases: []string{"u"},
			Value:   "17.0.0.0/24",
			Usage:   "The UE pool address",
		},
		&cli.StringFlag{
			Name:    "gnb-addr",
			Aliases: []string{"g"},
			Usage:   "The gNB address",
		},
		&cli.StringSliceFlag{
			Name:    "app-filter",
			Aliases: []string{"a"},
			Value:   []string{"ip:any:any:allow:100"},
			Usage:   "Specify an application filter. Format: '{ip | udp | tcp}:{IPv4 Prefix | any}:{<lower-L4-port>-<upper-L4-port> | any}:{allow | deny}:{rule-precedence}' . e.g. 'udp:10.0.0.0/8:80-88:allow:100'",
		},
		&cli.UintFlag{
			Name:    "qfi",
			Aliases: []string{"q"},
			Usage:   "The QFI value for QERs. Max value 64",
		},
	}
}

// validateCommonArgs validates common arguments
func validateCommonArgs(c *cli.Command) {
	baseID := c.Int("baseID")
	count := c.Int("count")

	if baseID <= 0 {
		logger.PfcpsimLog.Fatalln("baseID cannot be 0 or a negative number")
	}

	if count <= 0 {
		logger.PfcpsimLog.Fatalln("count cannot be 0 or a negative number")
	}
}

func GetSessionCommands() *cli.Command {
	return &cli.Command{
		Name:  "session",
		Usage: "Handle sessions",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create sessions",
				Flags: getCommonFlags(),
				Action: func(ctx context.Context, c *cli.Command) error {
					return sessionCreateAction(ctx, c)
				},
			},
			{
				Name:  "modify",
				Usage: "Modify sessions",
				Flags: append(getCommonFlags(), []cli.Flag{
					&cli.BoolFlag{
						Name:    "buffer",
						Aliases: []string{"b"},
						Usage:   "If set, downlink FARs will have the buffer flag set to true",
					},
					&cli.BoolFlag{
						Name:    "notifycp",
						Aliases: []string{"n"},
						Usage:   "Set true to have downlink FARs notify CP",
					},
				}...),
				Action: func(ctx context.Context, c *cli.Command) error {
					return sessionModifyAction(ctx, c)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete sessions",
				Flags: getCommonFlags(),
				Action: func(ctx context.Context, c *cli.Command) error {
					return sessionDeleteAction(ctx, c)
				},
			},
		},
	}
}

func sessionCreateAction(ctx context.Context, c *cli.Command) error {
	qfi := c.Uint("qfi")
	if qfi > 64 {
		logger.PfcpsimLog.Fatalf("qfi cannot be greater than 64. Provided qfi: %v", qfi)
	}

	client := connect()
	defer disconnect()

	validateCommonArgs(c)

	res, err := client.CreateSession(ctx, &pb.CreateSessionRequest{
		Count:         int32(c.Int("count")),
		BaseID:        int32(c.Int("baseID")),
		NodeBAddress:  c.String("gnb-addr"),
		UeAddressPool: c.String("ue-pool"),
		AppFilters:    c.StringSlice("app-filter"),
		Qfi:           int32(qfi),
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while creating sessions: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)
	return nil
}

func sessionModifyAction(ctx context.Context, c *cli.Command) error {
	client := connect()
	defer disconnect()

	validateCommonArgs(c)

	res, err := client.ModifySession(ctx, &pb.ModifySessionRequest{
		Count:         int32(c.Int("count")),
		BaseID:        int32(c.Int("baseID")),
		NodeBAddress:  c.String("gnb-addr"),
		UeAddressPool: c.String("ue-pool"),
		BufferFlag:    c.Bool("buffer"),
		NotifyCPFlag:  c.Bool("notifycp"),
		AppFilters:    c.StringSlice("app-filter"),
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while modifying sessions: %v", err)
	}

	logger.PfcpsimLog.Infof(res.Message)
	return nil
}

func sessionDeleteAction(ctx context.Context, c *cli.Command) error {
	client := connect()
	defer disconnect()

	validateCommonArgs(c)

	res, err := client.DeleteSession(ctx, &pb.DeleteSessionRequest{
		Count:  int32(c.Int("count")),
		BaseID: int32(c.Int("baseID")),
	})
	if err != nil {
		logger.PfcpsimLog.Fatalf("error while deleting sessions: %v", err)
	}

	logger.PfcpsimLog.Infoln(res.Message)
	return nil
}
