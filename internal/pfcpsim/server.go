/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package pfcpsim

import (
	"context"
	"fmt"
	"net"

	"github.com/c-robinson/iplib"
	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	log "github.com/sirupsen/logrus"
	ieLib "github.com/wmnsk/go-pfcp/ie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConcretePFCPSimServer implements the Protobuf interface and keeps a connection to a remote PFCP Agent peer.
// Its state is handled in internal/pfcpsim/state.go
type ConcretePFCPSimServer struct{}

func (P ConcretePFCPSimServer) Configure(ctx context.Context, request *pb.ConfigureRequest) (*pb.Response, error) {
	remotePeerAddress = request.RemotePeerAddress
	if net.ParseIP(remotePeerAddress) == nil {
		// Try to resolve hostname
		lookupHost, err := net.LookupHost(remotePeerAddress)
		if err != nil {
			errMsg := fmt.Sprintf("Could not retrieve hostname or address for remote peer: %s", remotePeerAddress)
			log.Error(errMsg)
			return &pb.Response{}, status.Error(codes.Aborted, errMsg)
		}
		remotePeerAddress = lookupHost[0]
	}

	if net.ParseIP(request.UpfN3Address) == nil {
		errMsg := fmt.Sprintf("Error while parsing UPF N3 address: %v", request.UpfN3Address)
		log.Error(errMsg)
		return &pb.Response{}, status.Error(codes.Aborted, errMsg)
	}

	upfN3Address = request.UpfN3Address

	if net.ParseIP(request.GnodeBAddress) == nil {
		log.Errorf("Could not retrieve IP address of gNodeB")
		return &pb.Response{}, status.Error(codes.Aborted, "Could not retrieve IP address of gNodeB")
	}

	gnodeBAddress = request.GnodeBAddress

	configurationMsg := fmt.Sprintf("Server is configured: \n\tRemote peer address: %v, N3 interface address: %v, gNodeB address: %v ", remotePeerAddress, upfN3Address, gnodeBAddress)
	log.Info(configurationMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    configurationMsg,
	}, nil
}

func (P ConcretePFCPSimServer) ConnectToRemotePeer(ctx context.Context, request *pb.EmptyRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	err := connectPFCPSim()
	if err != nil {
		return nil, err
	}
	infoMsg := "Connection to remote peer established"
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P ConcretePFCPSimServer) Interrupt(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	sim.DisconnectN4()

	remotePeerConnected = false

	infoMsg := "Connection to remote peer closed"
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P ConcretePFCPSimServer) Associate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	if !isConfigured() {
		log.Error("Server is not configured")
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	if !isRemotePeerConnected() {
		log.Error("PFCP agent is not connected to remote peer")
		return &pb.Response{}, status.Error(codes.Aborted, "PFCP agent is not connected to remote peer")
	}

	err := sim.SetupAssociation()
	if err != nil {
		log.Error(err.Error())
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	infoMsg := "Association established"
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P ConcretePFCPSimServer) Disassociate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	err := sim.TeardownAssociation()
	if err != nil {
		log.Error(err.Error())
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	infoMsg := "Association teardown completed"
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P ConcretePFCPSimServer) CreateSession(ctx context.Context, request *pb.CreateSessionRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	baseID := int(request.BaseID)
	count := int(request.Count)

	ueAddress := iplib.NextIP(net.IP(request.UeAddressPool)).String()

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
				WithN3Address(upfN3Address).
				WithSDFFilter("permit out ip from 0.0.0.0/0 to assigned").
				MarkAsUplink().
				BuildPDR(),

			// DownlinkPDR
			session.NewPDRBuilder().
				WithID(dowlinkPdrID).
				WithMethod(session.Create).
				WithPrecedence(100).
				WithUEAddress(ueAddress).
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
				WithDownlinkIP(gnodeBAddress).
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
			return &pb.Response{}, status.Error(codes.Internal, err.Error())
		}
		log.Infof("Saved session: %v", sess)
		insertSession(i, sess)
	}

	log.Infof("active sessions: %v, count: %v", len(activeSessions), count)

	infoMsg := fmt.Sprintf("%v sessions were established", count)
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P ConcretePFCPSimServer) ModifySession(ctx context.Context, request *pb.ModifySessionRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	// TODO handle buffer, notifyCP flags and 5G as well
	baseID := int(request.BaseID)
	count := int(request.Count)

	for i := baseID; i < (count + baseID); i++ {
		newFARs := []*ieLib.IE{
			// Downlink FAR
			session.NewFARBuilder().
				WithID(uint32(i)).
				WithMethod(session.Update).
				WithAction(session.ActionForward).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithTEID(uint32(i + 1)). // Same downlinkTEID that was generated in create sessions
				WithDownlinkIP(gnodeBAddress).
				BuildFAR(),
		}

		err := sim.ModifySession(getSession(i-baseID), nil, newFARs, nil)
		if err != nil {
			return &pb.Response{}, status.Error(codes.Internal, err.Error())
		}
	}

	log.Infof("active sessions: %v, count: %v", len(activeSessions), count)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    fmt.Sprintf("%v sessions correctly modified", count),
	}, nil
}

func (P ConcretePFCPSimServer) DeleteSession(ctx context.Context, request *pb.DeleteSessionRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	baseID := int(request.BaseID)
	count := int(request.Count)

	if len(activeSessions) < count {
		err := pfcpsim.NewNotEnoughSessionsError()
		log.Error(err.Error())
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	for i := baseID; i < (count + baseID); i++ {
		sess := getSession(i)
		log.Info("Got session: %v", sess)

		err := sim.DeleteSession(sess)
		if err != nil {
			log.Error(err.Error())
			return &pb.Response{}, status.Error(codes.Aborted, err.Error())
		}
		log.Infof("Session removed :%v", sess)
		// remove from activeSessions
		deleteSession(i)
	}

	log.Infof("active sessions: %v, count: %v", len(activeSessions), count)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    fmt.Sprintf("%v sessions deleted", count),
	}, nil
}