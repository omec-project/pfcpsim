// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package pfcpsim

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/c-robinson/iplib"
	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/logger"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	ieLib "github.com/wmnsk/go-pfcp/ie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// pfcpSimService implements the Protobuf interface and keeps a connection to a remote PFCP Agent peer.
// Its state is handled in internal/pfcpsim/state.go
type pfcpSimService struct{}

// SessionStep identifies the step in loops, used while creating/modifying/deleting sessions and rules IDs.
// It should be high enough to avoid IDs overlap when creating sessions. 5 Applications should be enough.
// In theory with ROC limitations, we should expect max 8 applications (5 explicit applications + 3 filters
// to deny traffic to the RFC1918 IPs, in case we have a ALLOW-PUBLIC)
const SessionStep = 10

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

func SetRemotePeer(addr string) {
	remotePeerAddress = addr
}

func SetUpfN3(addr string) {
	upfN3Address = addr
}

func (P pfcpSimService) Configure(ctx context.Context, request *pb.ConfigureRequest) (*pb.Response, error) {
	if net.ParseIP(request.UpfN3Address) == nil {
		errMsg := fmt.Sprintf("error while parsing UPF N3 address: %v", request.UpfN3Address)
		logger.PfcpsimLog.Errorln(errMsg)

		return &pb.Response{}, status.Error(codes.Aborted, errMsg)
	}
	// remotePeerAddress is validated in pfcpsim
	SetRemotePeer(request.RemotePeerAddress)
	SetUpfN3(request.UpfN3Address)

	configurationMsg := fmt.Sprintf(
		"Server is configured. Remote peer address: %v, N3 interface address: %v ",
		remotePeerAddress,
		upfN3Address,
	)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    configurationMsg,
	}, nil
}

func (P pfcpSimService) Associate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	if !isConfigured() {
		logger.PfcpsimLog.Errorln("server is not configured")
		return &pb.Response{}, status.Error(codes.Aborted, "Server is not configured")
	}

	if !isRemotePeerConnected() {
		if err := ConnectPFCPSim(); err != nil {
			errMsg := fmt.Sprintf("Could not connect to remote peer: %v", err)
			logger.PfcpsimLog.Error(errMsg)

			return &pb.Response{}, status.Error(codes.Aborted, errMsg)
		}
	}

	if err := sim.SetupAssociation(); err != nil {
		logger.PfcpsimLog.Errorln(err)
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	infoMsg := "Association established"
	logger.PfcpsimLog.Infoln(infoMsg)

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
		logger.PfcpsimLog.Errorln(err.Error())
		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	sim.DisconnectN4()

	remotePeerConnected = false

	infoMsg := "Association teardown completed and connection to remote peer closed"
	logger.PfcpsimLog.Infoln(infoMsg)

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
		errMsg := fmt.Sprintf("Could not parse Address Pool: %v", err)
		logger.PfcpsimLog.Errorln(errMsg)

		return &pb.Response{}, status.Error(codes.Aborted, errMsg)
	}

	var qfi uint8 = 0

	if request.Qfi != 0 {
		qfi = uint8(request.Qfi)
	}

	if err = isNumOfAppFiltersCorrect(request.AppFilters); err != nil {
		return &pb.Response{}, err
	}

	for i := baseID; i < (count*SessionStep + baseID); i = i + SessionStep {
		// using variables to ease comprehension on how rules are linked together
		uplinkTEID := uint32(i)

		ueAddress := iplib.NextIP(lastUEAddr)
		lastUEAddr = ueAddress

		sessQerID := uint32(0)

		var pdrs, fars, urrs []*ieLib.IE

		qers := []*ieLib.IE{
			// session QER
			session.NewQERBuilder().
				WithID(sessQerID).
				WithMethod(session.Create).
				WithUplinkMBR(60000).
				WithDownlinkMBR(60000).
				Build(),
		}

		// create as many PDRs, FARs and App QERs as the number of app filters provided through pfcpctl
		ID := uint16(i)

		for _, appFilter := range request.AppFilters {
			SDFFilter, gateStatus, precedence, err := ParseAppFilter(appFilter)
			if err != nil {
				return &pb.Response{}, status.Error(codes.Aborted, err.Error())
			}

			logger.PfcpsimLog.Infof("successfully parsed application filter. SDF Filter: %v", SDFFilter)

			uplinkPdrID := ID
			downlinkPdrID := ID + 1

			uplinkFarID := uint32(ID)
			downlinkFarID := uint32(ID + 1)

			uplinkAppQerID := uint32(ID)
			downlinkAppQerID := uint32(ID + 1)

			urrId := uint32(ID)
			urr := session.NewURRBuilder().
				WithID(urrId).
				WithMethod(session.Create).
				WithMeasurementMethod(0, 1, 0).
				WithMeasurementPeriod(1 * time.Second).
				WithReportingTrigger(session.ReportingTrigger{
					Flags: session.RPT_TRIG_PERIO,
				}).
				Build()

			urrs = append(urrs, urr)

			urr = session.NewURRBuilder().
				WithID(urrId+1).
				WithMethod(session.Create).
				WithMeasurementMethod(0, 1, 0).
				WithMeasurementPeriod(1*time.Second).
				WithReportingTrigger(session.ReportingTrigger{
					Flags: session.RPT_TRIG_VOLTH | session.RPT_TRIG_VOLQU,
				}).
				WithVolumeThreshold(7, 10000, 20000, 30000).
				WithVolumeQuota(7, 10000, 20000, 30000).
				Build()

			urrs = append(urrs, urr)

			uplinkPDR := session.NewPDRBuilder().
				WithID(uplinkPdrID).
				WithMethod(session.Create).
				WithTEID(uplinkTEID).
				WithFARID(uplinkFarID).
				AddQERID(sessQerID).
				AddQERID(uplinkAppQerID).
				WithN3Address(upfN3Address).
				WithSDFFilter(SDFFilter).
				WithPrecedence(precedence).
				MarkAsUplink().
				BuildPDR()

			downlinkPDR := session.NewPDRBuilder().
				WithID(downlinkPdrID).
				WithMethod(session.Create).
				WithPrecedence(precedence).
				WithUEAddress(ueAddress.String()).
				WithSDFFilter(SDFFilter).
				AddQERID(sessQerID).
				AddQERID(downlinkAppQerID).
				WithFARID(downlinkFarID).
				MarkAsDownlink().
				BuildPDR()

			pdrs = append(pdrs, uplinkPDR)
			pdrs = append(pdrs, downlinkPDR)

			uplinkFAR := session.NewFARBuilder().
				WithID(uplinkFarID).
				WithAction(session.ActionForward).
				WithDstInterface(ieLib.DstInterfaceCore).
				WithMethod(session.Create).
				BuildFAR()

			downlinkFAR := session.NewFARBuilder().
				WithID(downlinkFarID).
				WithAction(session.ActionDrop).
				WithMethod(session.Create).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithZeroBasedOuterHeaderCreation().
				BuildFAR()

			fars = append(fars, uplinkFAR)
			fars = append(fars, downlinkFAR)

			uplinkAppQER := session.NewQERBuilder().
				WithID(uplinkAppQerID).
				WithMethod(session.Create).
				WithQFI(qfi).
				WithUplinkMBR(50000).
				WithDownlinkMBR(30000).
				WithGateStatus(gateStatus).
				Build()

			downlinkAppQER := session.NewQERBuilder().
				WithID(downlinkAppQerID).
				WithMethod(session.Create).
				WithQFI(qfi).
				WithUplinkMBR(50000).
				WithDownlinkMBR(30000).
				WithGateStatus(gateStatus).
				Build()

			qers = append(qers, uplinkAppQER)
			qers = append(qers, downlinkAppQER)

			ID += 2
		}

		sess, err := sim.EstablishSession(pdrs, fars, qers, urrs)
		if err != nil {
			return &pb.Response{}, status.Error(codes.Internal, err.Error())
		}

		pfcpsim.InsertSession(i, sess)
	}

	infoMsg := fmt.Sprintf("%v sessions were established using %v as baseID", count, baseID)
	logger.PfcpsimLog.Infoln(infoMsg)

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

	if pfcpsim.GetActiveSessionNum() < count {
		err := pfcpsim.NewNotEnoughSessionsError()
		logger.PfcpsimLog.Errorln(err.Error())

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

	if err := isNumOfAppFiltersCorrect(request.AppFilters); err != nil {
		return &pb.Response{}, err
	}

	for i := baseID; i < (count*SessionStep + baseID); i = i + SessionStep {
		var newFARs []*ieLib.IE

		var newURRs []*ieLib.IE

		ID := uint32(i + 1)
		teid := uint32(i + 1)

		if request.BufferFlag || request.NotifyCPFlag {
			teid = 0 // When buffering, TEID = 0.
		}

		for range request.AppFilters {
			downlinkFAR := session.NewFARBuilder().
				WithID(ID). // Same FARID that was generated in create sessions
				WithMethod(session.Update).
				WithAction(actions).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithTEID(teid).
				WithDownlinkIP(nodeBaddress).
				BuildFAR()

			newFARs = append(newFARs, downlinkFAR)

			urrId := ID
			urr := session.NewURRBuilder().
				WithID(urrId).
				WithMethod(session.Update).
				WithMeasurementPeriod(1 * time.Second).
				Build()

			newURRs = append(newURRs, urr)

			ID += 2
		}

		sess, ok := pfcpsim.GetSession(i)
		if !ok {
			errMsg := fmt.Sprintf("Could not retrieve session with index %v", i)
			logger.PfcpsimLog.Errorln(errMsg)

			return &pb.Response{}, status.Error(codes.Internal, errMsg)
		}

		err := sim.ModifySession(sess, nil, newFARs, nil, newURRs)
		if err != nil {
			return &pb.Response{}, status.Error(codes.Internal, err.Error())
		}
	}

	infoMsg := fmt.Sprintf("%v sessions were modified", count)
	logger.PfcpsimLog.Infoln(infoMsg)

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

	if pfcpsim.GetActiveSessionNum() < count {
		err := pfcpsim.NewNotEnoughSessionsError()
		logger.PfcpsimLog.Error(err.Error())

		return &pb.Response{}, status.Error(codes.Aborted, err.Error())
	}

	for i := baseID; i < (count*SessionStep + baseID); i = i + SessionStep {
		sess, ok := pfcpsim.GetSession(i)
		if !ok {
			errMsg := "Session was nil. Check baseID"
			logger.PfcpsimLog.Errorln(errMsg)

			return &pb.Response{}, status.Error(codes.Aborted, errMsg)
		}

		err := sim.DeleteSession(sess)
		if err != nil {
			logger.PfcpsimLog.Errorln(err.Error())
			return &pb.Response{}, status.Error(codes.Aborted, err.Error())
		}
		// remove from activeSessions
		pfcpsim.RemoveSession(i)
	}

	infoMsg := fmt.Sprintf("%v sessions deleted; activeSessions: %v", count, pfcpsim.GetActiveSessionNum())
	logger.PfcpsimLog.Infoln(infoMsg)

	return &pb.Response{
		StatusCode: int32(codes.OK),
		Message:    infoMsg,
	}, nil
}
