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

const (
	defaultSDFfilter = "permit out ip from 0.0.0.0/0 to assigned"
)

// PFCPSimService implements the Protobuf interface and keeps a connection to a remote PFCP Agent peer.
// Its state is handled in internal/pfcpsim/state.go
type PFCPSimService struct{}

func (P PFCPSimService) Configure(ctx context.Context, request *pb.ConfigureRequest) (*pb.Response, error) {
	if net.ParseIP(request.UpfN3Address) == nil {
		errMsg := fmt.Sprintf("Error while parsing UPF N3 address: %v", request.UpfN3Address)
		log.Error(errMsg)
		return &pb.Response{}, status.Error(codes.Aborted, errMsg)
	}
	// remotePeerAddress is validated in pfcpsim
	remotePeerAddress = request.RemotePeerAddress
	upfN3Address = request.UpfN3Address

	configurationMsg := fmt.Sprintf("Server is configured. Remote peer address: %v, N3 interface address: %v ", remotePeerAddress, upfN3Address)
	log.Info(configurationMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    configurationMsg,
	}, nil
}

func (P PFCPSimService) Interrupt(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
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

func (P PFCPSimService) Associate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	if !isConfigured() {
		log.Error("Server is not configured")
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	if !isRemotePeerConnected() {
		if err := connectPFCPSim(); err != nil {
			errMsg := fmt.Sprintf("Could not connect to remote peer :%v", err)
			log.Error(errMsg)
			return &pb.Response{}, status.Error(codes.Aborted, errMsg)
		}
	}

	if err := sim.SetupAssociation(); err != nil {
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

func (P PFCPSimService) Disassociate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	if err := sim.TeardownAssociation(); err != nil {
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

func (P PFCPSimService) CreateSession(ctx context.Context, request *pb.CreateSessionRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	baseID := int(request.BaseID)
	count := int(request.Count)

	lastUEAddr, _, err := net.ParseCIDR(request.UeAddressPool)
	if err != nil {
		errMsg := fmt.Sprintf(" Could not parse Address Pool: %v", err)
		log.Error(errMsg)
		return &pb.Response{}, status.Error(codes.Aborted, errMsg)
	}

	nodebAddr := request.NodeBAddress

	for i := baseID; i < (count + baseID); i++ {
		// using variables to ease comprehension on how rules are linked together
		uplinkTEID := uint32(i)
		downlinkTEID := uint32(i + 1)

		ueAddress := iplib.NextIP(lastUEAddr)
		lastUEAddr = ueAddress

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
				WithSDFFilter(defaultSDFfilter).
				MarkAsUplink().
				BuildPDR(),

			// DownlinkPDR
			session.NewPDRBuilder().
				WithID(dowlinkPdrID).
				WithMethod(session.Create).
				WithPrecedence(100).
				WithUEAddress(ueAddress.String()).
				WithSDFFilter(defaultSDFfilter).
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
				WithDownlinkIP(nodebAddr).
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
		insertSession(i, sess)
	}

	infoMsg := fmt.Sprintf("%v sessions were established using %v as baseID ", count, baseID)
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P PFCPSimService) ModifySession(ctx context.Context, request *pb.ModifySessionRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	// TODO handle buffer, notifyCP flags and 5G as well
	baseID := int(request.BaseID)
	count := int(request.Count)
	nodeBaddress := request.NodeBAddress

	if len(activeSessions) < count {
		err := pfcpsim.NewNotEnoughSessionsError()
		log.Error(err)
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	for i := baseID; i < (count + baseID); i++ {
		newFARs := []*ieLib.IE{
			// Downlink FAR
			session.NewFARBuilder().
				WithID(uint32(i + 1)). // Same FARID that was generated in create sessions
				WithMethod(session.Update).
				WithAction(session.ActionForward).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithTEID(uint32(i + 1)). // Same downlinkTEID that was generated in create sessions
				WithDownlinkIP(nodeBaddress).
				BuildFAR(),
		}

		sess, ok := getSession(i)
		if !ok {
			errMsg := fmt.Sprintf("Could not retrieve session with index %v", i)
			log.Error(errMsg)
			return &pb.Response{}, status.Error(codes.Internal, errMsg)
		}

		err := sim.ModifySession(sess, nil, newFARs, nil)
		if err != nil {
			return &pb.Response{}, status.Error(codes.Internal, err.Error())
		}
	}

	infoMsg := fmt.Sprintf("%v sessions were modified", count)
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P PFCPSimService) DeleteSession(ctx context.Context, request *pb.DeleteSessionRequest) (*pb.Response, error) {
	if !isConfigured() {
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	baseID := int(request.BaseID)
	count := int(request.Count)

	if len(activeSessions) < count {
		err := pfcpsim.NewNotEnoughSessionsError()
		log.Error(err)
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	for i := baseID; i < (count + baseID); i++ {
		sess, ok := getSession(i)
		if !ok {
			errMsg := "Session was nil. Check baseID"
			log.Error(errMsg)
			return &pb.Response{}, status.Error(codes.Aborted, errMsg)
		}

		err := sim.DeleteSession(sess)
		if err != nil {
			log.Error(err.Error())
			return &pb.Response{}, status.Error(codes.Aborted, err.Error())
		}
		// remove from activeSessions
		deleteSession(i)
	}

	infoMsg := fmt.Sprintf("%v sessions deleted; activeSessions: %v", count, len(activeSessions))
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}
