/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"net"
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
)

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
