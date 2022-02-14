/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package pfcpsim

import (
	"net"

	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
)

func connectPFCPSim() error {
	if sim == nil {
		localAddr, err := getLocalAddress(interfaceName)
		if err != nil {
			return err
		}

		sim = pfcpsim.NewPFCPClient(localAddr.String())
	}

	err := sim.ConnectN4(remotePeerAddress)
	if err != nil {
		return err
	}

	remotePeerConnected = true

	return nil
}

func isConfigured() bool {
	if upfN3Address != "" && remotePeerAddress != "" {
		return true
	}

	return false
}

func isRemotePeerConnected() bool {
	return remotePeerConnected
}

// getLocalAddress returns the first IP address of the interfaceName, if specified,
// otherwise returns the IP address of the first non-loopback interface
// Returns error if fail occurs at any stage.
func getLocalAddress(interfaceName string) (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	if interfaceName != "" {
		// Interface name is specified. Use it.
		interfaceAddrs, err := net.InterfaceByName(interfaceName)
		if err != nil {
			return nil, err
		}

		addrs, _ = interfaceAddrs.Addrs()
	}

	for _, address := range addrs {
		// Check address type to be non-loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, pfcpsim.NewNoValidInterfaceError()
}
