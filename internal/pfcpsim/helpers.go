/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package pfcpsim

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"

	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
)

func connectPFCPSim() error {
	if sim == nil {
		localAddress, err := getLocalAddress()
		if err != nil {
			return err
		}

		sim = pfcpsim.NewPFCPClient(localAddress)
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

// getLocalAddress retrieves local address in string format to use when establishing a connection with PFCP agent
func getLocalAddress() (string, error) {
	// cmd to run for darwin platforms
	cmd := "route -n get default | grep 'interface:' | grep -o '[^ ]*$'"

	if runtime.GOOS != "darwin" {
		// assuming linux platform
		cmd = "route | grep '^default' | grep -o '[^ ]*$'"
	}

	cmdOutput, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
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
		return ip.String(), nil
	}

	return "", fmt.Errorf("could not find interface: %v", interfaceName)
}
