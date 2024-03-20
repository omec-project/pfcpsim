// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package config

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	defaultgRPCServerAddress = "localhost:54321"
)

var GlobalOptions struct {
	Server string `short:"s" long:"server" default:"" value-name:"SERVER:PORT" description:"IP/Host and port of pfcpsim gRPC Server"`
}

type GlobalConfigSpec struct {
	Server string
}

var GlobalConfig = GlobalConfigSpec{
	Server: defaultgRPCServerAddress,
}

func ProcessGlobalOptions() {
	// Override from environment
	serverFromEnv, present := os.LookupEnv("PFCPSIM_SERVER")
	if present {
		GlobalConfig.Server = serverFromEnv
	}

	// Override from command line
	if GlobalOptions.Server != "" {
		GlobalConfig.Server = GlobalOptions.Server
	}

	// Generate error messages for required settings
	if GlobalConfig.Server == "" {
		log.Fatal("Server is not set. Please use the -s option")
	}

	// Try to resolve hostname if provided for the server
	if host, port, err := net.SplitHostPort(GlobalConfig.Server); err == nil {
		if addrs, err := net.LookupHost(host); err == nil {
			GlobalConfig.Server = net.JoinHostPort(addrs[0], port)
		}
	}

	log.Debugf("ServerAddress: %v", GlobalOptions.Server)
}
