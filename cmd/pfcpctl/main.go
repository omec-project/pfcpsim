// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/omec-project/pfcpsim/internal/pfcpctl/commands"
	"github.com/omec-project/pfcpsim/internal/pfcpctl/config"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{}
	app.Name = path.Base(os.Args[0])
	app.Usage = "PFCP Control CLI"
	app.Flags = config.GetGlobalFlags()
	app.Before = beforeAction
	app.Commands = []*cli.Command{
		// Service commands
		commands.GetServiceCommands(),
		// Session commands
		commands.GetSessionCommands(),
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}

func beforeAction(ctx context.Context, c *cli.Command) (context.Context, error) {
	// Set global options from flags
	config.SetGlobalOptionsFromCli(c)
	// Set server address and configure other parameters
	config.ProcessGlobalOptions()
	return ctx, nil
}
