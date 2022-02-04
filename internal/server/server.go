/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	ieLib "github.com/wmnsk/go-pfcp/ie"
)

// pfcpSimServer implements the Protobuf methods and keeps a connection to a remote PFCP Agent peer.
// It stores only 'static' values. Its state is handled in internal/server/state.go
type pfcpSimServer struct {
	upfAddress    string
	nodeBAddress  string
	ueAddressPool string
}

// getLocalAddress retrieves local address to use when establishing a connection with PFCP agent
func getLocalAddress() (net.IP, error) {
	// cmd to run for darwin platforms
	cmd := "route -n get default | grep 'interface:' | grep -o '[^ ]*$'"

	if runtime.GOOS != "darwin" {
		// assuming linux platform
		cmd = "route | grep '^default' | grep -o '[^ ]*$'"
	}

	cmdOutput, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
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
		return ip, nil
	}

	return nil, fmt.Errorf("could not find interface: %v", interfaceName)
}

func NewPFCPSimServer(remotePeerAddr string, upfAddress string, nodeBAddress string, ueAddressPool string) (*pfcpSimServer, error) {
	lAddr, err := getLocalAddress()
	if err != nil {
		return nil, err
	}

	// Connect internal pfcpSim to remote Peer
	pfcpSim = pfcpsim.NewPFCPClient(lAddr.String())
	err = pfcpSim.ConnectN4(remotePeerAddr)
	if err != nil {
		return nil, err
	}

	return &pfcpSimServer{
		upfAddress:    upfAddress,
		nodeBAddress:  nodeBAddress,
		ueAddressPool: ueAddressPool,
	}, nil
}

func (P pfcpSimServer) Associate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	err := pfcpSim.SetupAssociation()
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Association completed",
	}, nil
}

func (P pfcpSimServer) Disassociate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	err := pfcpSim.TeardownAssociation()
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
	count := int(request.Count) // cast int32 to int

	for i := baseID; i < (count + baseID); i++ {
		// using variables to ease comprehension on how rules are linked together
		uplinkTEID := uint32(i + 10)
		downlinkTEID := uint32(i + 11)

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
				WithN3Address(P.upfAddress).
				WithSDFFilter("permit out ip from any to assigned").
				MarkAsUplink().
				BuildPDR(),

			// DownlinkPDR
			session.NewPDRBuilder().
				WithID(dowlinkPdrID).
				WithMethod(session.Create).
				WithPrecedence(100).
				WithUEAddress(getNextUEAddress(P.ueAddressPool).String()).
				WithSDFFilter("permit out ip from any to assigned").
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
				WithDownlinkIP(P.nodeBAddress).
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

		sess, err := pfcpSim.EstablishSession(pdrs, fars, qers)
		if err != nil {
			return nil, err
		}

		ctx := &pfcpClientContext{
			session:      sess,
			pdrs:         pdrs,
			fars:         fars,
			qers:         qers,
			downlinkTEID: downlinkTEID,
		}

		addSessionContext(ctx)
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
					WithDownlinkIP(P.nodeBAddress).
					BuildFAR(),
			}

			err = pfcpSim.ModifySession(ctx.session, nil, newFARs, nil)
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
