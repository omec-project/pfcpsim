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

	"github.com/c-robinson/iplib"
	pb "github.com/omec-project/pfcpsim/api"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	log "github.com/sirupsen/logrus"
	ieLib "github.com/wmnsk/go-pfcp/ie"
)

type pfcpClientContext struct {
	session *pfcpsim.PFCPSession

	pdrs []*ieLib.IE
	fars []*ieLib.IE
	qers []*ieLib.IE
}

// pfcpSimServer implements the Protobuf methods and keeps a connection to a remote PFCP Agent peer.
type pfcpSimServer struct {
	// Emulates 5G SMF/ 4G SGW
	client *pfcpsim.PFCPClient

	activeSessions []*pfcpClientContext

	upfAddress    string
	nodeBAddress  string
	ueAddressPool string

	lastUEAddress net.IP
}

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

	cl := pfcpsim.NewPFCPClient(lAddr.String())
	err = cl.ConnectN4(remotePeerAddr)
	if err != nil {
		return nil, err
	}

	return &pfcpSimServer{
		client:         cl,
		upfAddress:     upfAddress,
		nodeBAddress:   nodeBAddress,
		ueAddressPool:  ueAddressPool,
		activeSessions: make([]*pfcpClientContext, 0),
	}, nil
}

// getNextUEAddress retrieves the next available IP address from ueAddressPool
func (P *pfcpSimServer) getNextUEAddress() net.IP {
	if P.lastUEAddress != nil {
		P.lastUEAddress = iplib.NextIP(P.lastUEAddress)
		return P.lastUEAddress
	}

	// TODO handle case net IP is full
	ueIpFromPool, _, _ := net.ParseCIDR(P.ueAddressPool)
	P.lastUEAddress = iplib.NextIP(ueIpFromPool)
	return P.lastUEAddress
}

// createSessions create 'count' sessions incrementally.
// Once created, the sessions are established through PFCP client.
func (P *pfcpSimServer) createSessions(count int) {

}

func (P pfcpSimServer) SetLogLevel(ctx context.Context, level *pb.LogLevel) (*pb.LogLevel, error) {
	//TODO implement me
	panic("implement me")
}

func (P pfcpSimServer) Associate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	err := P.client.SetupAssociation()
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Association completed",
	}, nil
}

func (P pfcpSimServer) Disassociate(ctx context.Context, empty *pb.EmptyRequest) (*pb.Response, error) {
	err := P.client.TeardownAssociation()
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Association teardown completed",
	}, nil
}

func (P pfcpSimServer) CreateSession(ctx context.Context, request *pb.CreateSessionRequest) (*pb.Response, error) {
	baseID := len(P.activeSessions) + 1
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
				WithUEAddress(P.getNextUEAddress().String()).
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

		sess, err := P.client.EstablishSession(pdrs, fars, qers)
		if err != nil {
			log.Errorf("Error while establishing sessions: %v", err)
			return nil, err
		}

		P.activeSessions = append(P.activeSessions, &pfcpClientContext{
			session: sess,
			pdrs:    pdrs,
			fars:    fars,
			qers:    qers,
		},
		)

		log.Infof("Created new PFCP session")
	}

	return &pb.Response{
		StatusCode: http.StatusOK,
		Message:    "Created Session",
	}, nil
}

func (P pfcpSimServer) ModifySession(ctx context.Context, request *pb.ModifySessionRequest) (*pb.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (P pfcpSimServer) DeleteSession(ctx context.Context, request *pb.DeleteSessionRequest) (*pb.Response, error) {
	//TODO implement me
	panic("implement me")
}
