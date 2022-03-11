// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

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

// pfcpSimService implements the Protobuf interface and keeps a connection to a remote PFCP Agent peer.
// Its state is handled in internal/pfcpsim/state.go
type pfcpSimService struct{}

func NewPFCPSimService(iface string) *pfcpSimService {
	interfaceName = iface
	return &pfcpSimService{}
}

func checkServerStatus() error {
	if !isConfigured() {
		return status.Error(codes.Aborted, "Server is not configured")
	}

	if !isRemotePeerConnected() {
		return status.Error(codes.Aborted, "Server is not associated")
	}

	return nil
}

func (P pfcpSimService) Configure(ctx context.Context, request *pb.ConfigureRequest) (*pb.Response, error) {
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

func (P pfcpSimService) Associate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
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

func (P pfcpSimService) Disassociate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	if err := checkServerStatus(); err != nil {
		return &pb.Response{}, err
	}

	if err := sim.TeardownAssociation(); err != nil {
		log.Error(err.Error())
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	sim.DisconnectN4()

	remotePeerConnected = false

	infoMsg := "Association teardown completed and connection to remote peer closed"
	log.Info(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}

func (P pfcpSimService) CreateSession(ctx context.Context, request *pb.CreateSessionRequest) (*pb.Response, error) {
	if err := checkServerStatus(); err != nil {
		return &pb.Response{}, err
	}

	baseID := int(request.BaseID)
	count := int(request.Count)

	lastUEAddr, _, err := net.ParseCIDR(request.UeAddressPool)
	if err != nil {
		errMsg := fmt.Sprintf(" Could not parse Address Pool: %v", err)
		log.Error(errMsg)
		return &pb.Response{}, status.Error(codes.Aborted, errMsg)
	}

	var SDFFilter = ""
	var qfi, gateStatus uint8 = 0, ieLib.GateStatusOpen

	if request.GateClosed {
		gateStatus = ieLib.GateStatusClosed
	}

	if request.Qfi != 0 {
		qfi = uint8(request.Qfi)
	}

	if request.SdfFilter != "" {
		SDFFilter = request.SdfFilter
	}

	for i := baseID; i < (count*2 + baseID); i = i + 2 {
		// using variables to ease comprehension on how rules are linked together
		uplinkTEID := uint32(i)

		ueAddress := iplib.NextIP(lastUEAddr)
		lastUEAddr = ueAddress

		uplinkFarID := uint32(i)
		downlinkFarID := uint32(i + 1)

		uplinkPdrID := uint16(i)
		dowlinkPdrID := uint16(i + 1)

		sessQerID := uint32(i + 3)

		uplinkAppQerID := uint32(i)
		downlinkAppQerID := uint32(i + 1)

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
				WithSDFFilter(SDFFilter).
				WithPrecedence(100).
				MarkAsUplink().
				BuildPDR(),

			// DownlinkPDR
			session.NewPDRBuilder().
				WithID(dowlinkPdrID).
				WithMethod(session.Create).
				WithPrecedence(100).
				WithUEAddress(ueAddress.String()).
				WithSDFFilter(SDFFilter).
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
				WithZeroBasedOuterHeaderCreation().
				BuildFAR(),
		}

		qers := []*ieLib.IE{
			// TODO make rates configurable by pfcpctl
			// session QER
			session.NewQERBuilder().
				WithID(sessQerID).
				WithMethod(session.Create).
				WithUplinkMBR(60000).
				WithDownlinkMBR(60000).
				Build(),

			// Uplink application QER
			session.NewQERBuilder().
				WithID(uplinkAppQerID).
				WithMethod(session.Create).
				WithQFI(qfi).
				WithUplinkMBR(50000).
				WithDownlinkMBR(30000).
				WithGateStatus(gateStatus).
				Build(),

			// Downlink application QER
			session.NewQERBuilder().
				WithID(downlinkAppQerID).
				WithMethod(session.Create).
				WithQFI(qfi).
				WithUplinkMBR(50000).
				WithDownlinkMBR(30000).
				WithGateStatus(gateStatus).
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

func (P pfcpSimService) ModifySession(ctx context.Context, request *pb.ModifySessionRequest) (*pb.Response, error) {
	if err := checkServerStatus(); err != nil {
		return &pb.Response{}, err
	}

	// TODO add 5G mode
	baseID := int(request.BaseID)
	count := int(request.Count)
	nodeBaddress := request.NodeBAddress

	if len(activeSessions) < count {
		err := pfcpsim.NewNotEnoughSessionsError()
		log.Error(err)
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	var actions uint8 = 0

	if request.BufferFlag || request.NotifyCPFlag {
		// We currently support only both flags set
		actions |= session.ActionNotify
		actions |= session.ActionBuffer
	} else {
		// If no flag was passed, default action is Forward
		actions |= session.ActionForward
	}

	for i := baseID; i < (count*2 + baseID); i = i + 2 {
		teid := uint32(i + 1)

		if request.BufferFlag || request.NotifyCPFlag {
			teid = 0 // When buffering, TEID = 0.
		}

		newFARs := []*ieLib.IE{
			// Downlink FAR
			session.NewFARBuilder().
				WithID(uint32(i + 1)). // Same FARID that was generated in create sessions
				WithMethod(session.Update).
				WithAction(actions).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithTEID(teid).
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

func (P pfcpSimService) DeleteSession(ctx context.Context, request *pb.DeleteSessionRequest) (*pb.Response, error) {
	if err := checkServerStatus(); err != nil {
		return &pb.Response{}, err
	}

	baseID := int(request.BaseID)
	count := int(request.Count)

	if len(activeSessions) < count {
		err := pfcpsim.NewNotEnoughSessionsError()
		log.Error(err)
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	for i := baseID; i < (count*2 + baseID); i = i + 2 {
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
