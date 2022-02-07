/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	ieLib "github.com/wmnsk/go-pfcp/ie"
)

// pfcpSimServer implements the Protobuf methods and keeps a connection to a remote PFCP Agent peer.
// Its state is handled in internal/server/state.go
type pfcpSimServer struct{}

func (P pfcpSimServer) Recover(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	err := connectPFCPSim()
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Emulating Crash",
	}, nil
}

func (P pfcpSimServer) Interrupt(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	sim.DisconnectN4()

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Emulating Crash",
	}, nil
}

func NewPFCPSimServer(remotePeerAddr string, upfAddr string, nodeBAddr string, ueAddrPool string) (*pfcpSimServer, error) {
	var err error
	localAddress, err = getLocalAddress()
	if err != nil {
		return nil, err
	}

	remotePeerAddress = remotePeerAddr
	upfAddress = upfAddr
	nodeBAddress = nodeBAddr
	ueAddressPool = ueAddrPool

	// Connect internal PFCPSim to remote Peer
	err = connectPFCPSim()
	if err != nil {
		return nil, err
	}

	return &pfcpSimServer{}, nil
}

func (P pfcpSimServer) Associate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	err := sim.SetupAssociation()
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Association completed",
	}, nil
}

func (P pfcpSimServer) Disassociate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	err := sim.TeardownAssociation()
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Association teardown completed",
	}, nil
}

func (P pfcpSimServer) CreateSession(ctx context.Context, request *pb.CreateSessionRequest) (*pb.Response, error) {
	sessions := getActiveSessions()

	baseID := len(*sessions) + 1
	count := int(request.Count)

	for i := baseID; i < (count + baseID); i++ {
		// using variables to ease comprehension on how rules are linked together
		uplinkTEID := uint32(i)
		downlinkTEID := uint32(i + 1)

		uplinkFarID := uint32(i)
		downlinkFarID := uint32(i + 1)

		uplinkPdrID := uint16(i)
		dowlinkPdrID := uint16(i + 1)

		sessQerID := uint32(i + 3)

		appQerID := uint32(i)

		uplinkAppQerID := appQerID
		downlinkAppQerID := appQerID + 1

		pdrs := []*ieLib.IE{
			// UplinkPDR
			session.NewPDRBuilder().
				WithID(uplinkPdrID).
				WithMethod(session.Create).
				WithTEID(uplinkTEID).
				WithFARID(uplinkFarID).
				AddQERID(sessQerID).
				AddQERID(uplinkAppQerID).
				WithN3Address(upfAddress).
				WithSDFFilter("permit out ip from 0.0.0.0/0 to assigned").
				MarkAsUplink().
				BuildPDR(),

			// DownlinkPDR
			session.NewPDRBuilder().
				WithID(dowlinkPdrID).
				WithMethod(session.Create).
				WithPrecedence(100).
				WithUEAddress(getNextUEAddress(ueAddressPool).String()).
				WithSDFFilter("permit out ip from 0.0.0.0/0 to assigned").
				AddQERID(sessQerID).
				AddQERID(downlinkAppQerID).
				WithFARID(downlinkFarID).
				MarkAsDownlink().
				BuildPDR(),
		}

		fars := []*ieLib.IE{
			// UplinkFAR
			session.NewFARBuilder().
				WithID(uplinkFarID).
				WithAction(session.ActionForward).
				WithDstInterface(ieLib.DstInterfaceCore).
				WithMethod(session.Create).
				BuildFAR(),

			// DownlinkFAR
			session.NewFARBuilder().
				WithID(downlinkFarID).
				WithAction(session.ActionDrop).
				WithMethod(session.Create).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithTEID(downlinkTEID).
				WithDownlinkIP(nodeBAddress).
				BuildFAR(),
		}

		qers := []*ieLib.IE{
			// session QER
			session.NewQERBuilder().
				WithID(sessQerID).
				WithMethod(session.Create).
				WithQFI(0x09).
				WithUplinkMBR(50000).
				WithDownlinkMBR(50000).
				Build(),

			// application QER
			session.NewQERBuilder().
				WithID(appQerID).
				WithMethod(session.Create).
				WithQFI(0x08).
				WithUplinkMBR(50000).
				WithUplinkGBR(50000).
				WithDownlinkMBR(30000).
				WithUplinkGBR(30000).
				Build(),
		}

		sess, err := sim.EstablishSession(pdrs, fars, qers)
		if err != nil {
			return nil, err
		}

		addSessionContext(&pfcpClientContext{
			session:      sess,
			pdrs:         pdrs,
			fars:         fars,
			qers:         qers,
			downlinkTEID: downlinkTEID,
		})
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("%v sessions were established", count),
	}, nil
}

func (P pfcpSimServer) ModifySession(ctx context.Context, request *pb.ModifySessionRequest) (*pb.Response, error) {
	sessions := getActiveSessions()

	count := int(request.Count)

	if len(*sessions) < count {
		return nil, pfcpsim.NewNotEnoughSessionsError()
	}

	for i, ctx := range *sessions {
		if i >= count {
			// Modify only 'count' sessions
			break
		}

		for _, far := range ctx.fars {
			action, err := far.ApplyAction()
			if err != nil {
				return nil, err
			}

			if !(action == session.ActionDrop) {
				// Updating only FARs with ActionDrop.
				continue
			}

			oldFarID, err := far.FARID()
			if err != nil {
				return nil, err
			}

			// TODO handle buffer and notifyCP flags and 5G as well
			newFARs := []*ieLib.IE{
				// Downlink FAR
				session.NewFARBuilder().
					WithID(oldFarID).
					WithMethod(session.Update).
					WithAction(session.ActionForward).
					WithDstInterface(ieLib.DstInterfaceCore).
					WithTEID(ctx.downlinkTEID).
					WithDownlinkIP(nodeBAddress).
					BuildFAR(),
			}

			err = sim.ModifySession(ctx.session, nil, newFARs, nil)
			if err != nil {
				return nil, err
			}

			ctx.fars = append(ctx.fars, newFARs...) // save new sent FAR // FIXME should old FARs be replaced?
		}
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("%v sessions correctly modified", count),
	}, nil
}

func (P pfcpSimServer) DeleteSession(ctx context.Context, request *pb.DeleteSessionRequest) (*pb.Response, error) {
	sessions := getActiveSessions()

	count := int(request.Count)

	if len(*sessions) < count {
		return nil, pfcpsim.NewNotEnoughSessionsError()
	}

	for i := count; i > 0; i-- {
		deleteSessionContext()
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("%v sessions deleted", count),
	}, nil
}
