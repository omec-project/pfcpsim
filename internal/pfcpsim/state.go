/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package pfcpsim

import (
	"sync"

	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
)

var (
	activeSessions     = make(map[int]*pfcpsim.PFCPSession, 0)
	lockActiveSessions = new(sync.Mutex)

	remotePeerAddress string
	upfN3Address      string

	interfaceName string
	pcapPath string
	snifferDoneChannel chan bool
	isSnifferStarted bool

	// Emulates 5G SMF/ 4G SGW
	sim                 *pfcpsim.PFCPClient
	remotePeerConnected bool
)

func startSniffer() {
	if isSnifferStarted {
		return
	}

	snifferDoneChannel = make(chan bool)

	go func() {
		err := sniffer(snifferDoneChannel)
		if err != nil {
			return
		}
	}()

	isSnifferStarted = true
}

func stopSniffer() {
	if !isSnifferStarted {
		return
	}

	select {
	case snifferDoneChannel <- true:
		isSnifferStarted = false
		close(snifferDoneChannel)
		return
	default:
		// Sniffer was not started
		isSnifferStarted = false
		close(snifferDoneChannel)
	}
}

func insertSession(index int, session *pfcpsim.PFCPSession) {
	lockActiveSessions.Lock()
	defer lockActiveSessions.Unlock()

	activeSessions[index] = session
}

func getSession(index int) (*pfcpsim.PFCPSession, bool) {
	element, ok := activeSessions[index]
	return element, ok
}

func deleteSession(index int) {
	lockActiveSessions.Lock()
	defer lockActiveSessions.Unlock()

	delete(activeSessions, index)
}
