// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package config

import (
	"net"
	"os"

	"github.com/omec-project/pfcpsim/logger"
	"github.com/urfave/cli/v3"
)

const (
	defaultgRPCServerAddress = "localhost:54321"
)

type GlobalConfigSpec struct {
	Server string
}

var GlobalConfig = GlobalConfigSpec{
	Server: defaultgRPCServerAddress,
}

// GetGlobalFlags returns the global flags for urfave/cli
func GetGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "server",
			Aliases: []string{"s"},
			Value:   "",
			Usage:   "gRPC Server IP/Host and port (SERVER:PORT)",
		},
	}
}

// SetGlobalOptionsFromCli sets global options from CLI context
func SetGlobalOptionsFromCli(c *cli.Command) {
	serverFlag := c.String("server")

	// Start with default
	GlobalConfig.Server = defaultgRPCServerAddress

	// Override from environment
	if serverFromEnv, present := os.LookupEnv("PFCPSIM_SERVER"); present {
		GlobalConfig.Server = serverFromEnv
	}

	// Override from command line (highest priority)
	if serverFlag != "" {
		GlobalConfig.Server = serverFlag
	}
}

func ProcessGlobalOptions() {
	// Generate error messages for required settings
	if GlobalConfig.Server == "" {
		logger.PfcpsimLog.Fatalln("server is not set. Please use the -s option or set PFCPSIM_SERVER environment variable")
	}

	// Try to resolve hostname if provided for the server
	if host, port, err := net.SplitHostPort(GlobalConfig.Server); err == nil {
		if addrs, err := net.LookupHost(host); err == nil {
			GlobalConfig.Server = net.JoinHostPort(addrs[0], port)
		}
	}

	logger.PfcpsimLog.Debugf("serverAddress: %v", GlobalConfig.Server)
}
