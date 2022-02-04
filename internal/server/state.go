/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/c-robinson/iplib"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	ieLib "github.com/wmnsk/go-pfcp/ie"
)

type pfcpClientContext struct {
	session *pfcpsim.PFCPSession

	pdrs []*ieLib.IE
	fars []*ieLib.IE
	qers []*ieLib.IE
	// Needed when updating sessions.
	downlinkTEID uint32
}

var (
	activeSessions = make([]*pfcpClientContext, 0)
	sessionsLock   = new(sync.Mutex)

	// Keeps track of 'leased' IPs to UEs from ip pool
	lastUEAddress net.IP
	addrLock      = new(sync.Mutex)

	localAddress      string
	remotePeerAddress string
	upfAddress        string
	nodeBAddress      string
	ueAddressPool     string

	// Emulates 5G SMF/ 4G SGW
	sim *pfcpsim.PFCPClient
)

func connectPFCPSim() error {
	if sim == nil {
		sim = pfcpsim.NewPFCPClient(localAddress)
	}

	err := sim.ConnectN4(remotePeerAddress)
	if err != nil {
		return err
	}

	return nil
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

func getActiveSessions() *[]*pfcpClientContext {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	return &activeSessions
}

func deleteSessionContext() {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	if len(activeSessions) > 0 {
		// pop first element
		activeSessions = activeSessions[:len(activeSessions)-1]
	}
}

func addSessionContext(sessionContext *pfcpClientContext) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	activeSessions = append(activeSessions, sessionContext)
}

// getNextUEAddress retrieves the next available IP address from ueAddressPool
func getNextUEAddress(addressPool string) net.IP {
	addrLock.Lock()
	defer addrLock.Unlock()

	if lastUEAddress != nil {
		lastUEAddress = iplib.NextIP(lastUEAddress)
		return lastUEAddress
	}

	// TODO handle case net IP is full
	ueIpFromPool, _, _ := net.ParseCIDR(addressPool)
	lastUEAddress = iplib.NextIP(ueIpFromPool)
	return lastUEAddress
}
