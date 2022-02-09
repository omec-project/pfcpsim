/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package pfcpsim

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	ieLib "github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"
)

func Test_ListenN4(t *testing.T) {
	wg := new(sync.WaitGroup)
	CPSim := NewPFCPClient("127.0.0.1")

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		// Start CP in goroutine, awaiting UP to start association
		defer wg.Done()

		err := CPSim.ListenN4()
		require.NoError(t, err)
		// Once association is set up, ListenN4 returns
		// Now let CP try to establish a fake session
		_, err = CPSim.EstablishSession(nil, nil, nil)
		require.NoError(t, err)
	}(wg)

	// wait for CPSim to start listening
	time.Sleep(time.Millisecond * 50)

	// Start mockUPF
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		mockUPFConn, err := net.Dial("udp", fmt.Sprintf(":%v", PFCPStandardPort))
		require.NoError(t, err)

		assocReq := message.NewAssociationSetupRequest(
			1,
			ieLib.NewRecoveryTimeStamp(time.Now()),
			ieLib.NewNodeID("127.0.0.1", "", ""),
		)

		marshal, _ := assocReq.Marshal()

		_, err = mockUPFConn.Write(marshal)
		require.NoError(t, err)

		buf := make([]byte, 1500)
		// Keep reading until a session establish request is received
		for {
			n, err := mockUPFConn.Read(buf)
			require.NoError(t, err)

			msg, err := message.Parse(buf[:n])
			require.NoError(t, err)

			if _, ok := msg.(*message.SessionEstablishmentRequest); ok {
				mockResponse := message.NewSessionEstablishmentResponse(
					0,
					0,
					999,
					1,
					0,
					ieLib.NewFSEID(998, net.ParseIP("127.0.0.1"), nil),
					ieLib.NewCause(ieLib.CauseRequestAccepted),
				)

				b, _ := mockResponse.Marshal()
				_, err = mockUPFConn.Write(b)
				require.NoError(t, err)

				return
			}
		}
	}(wg)

	wg.Wait()
}
