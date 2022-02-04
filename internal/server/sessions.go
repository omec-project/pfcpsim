/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"sync"

	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	log "github.com/sirupsen/logrus"
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
	lock           = new(sync.Mutex)
)

func getActiveSessions() *[]*pfcpClientContext {
	log.Infof(" Sessions: %v", activeSessions)
	return &activeSessions
}

func deleteSessionContext() {
	lock.Lock()
	defer lock.Unlock()
	if len(activeSessions) > 0 {
		// pop first element
		activeSessions = activeSessions[:len(activeSessions)-1]
	}
}

func addSessionContext(sessionContext *pfcpClientContext) {
	lock.Lock()
	defer lock.Unlock()

	activeSessions = append(activeSessions, sessionContext)
}
