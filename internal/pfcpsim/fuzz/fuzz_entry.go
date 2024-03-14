package fuzz

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/omec-project/pfcpsim/internal/pfcpsim"
	sim "github.com/omec-project/pfcpsim/pkg/pfcpsim"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	log "github.com/sirupsen/logrus"
	ieLib "github.com/wmnsk/go-pfcp/ie"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SVC_INIT = iota
	SVC_CONNECTED
	SVC_ASSOCIATED
)

var (
	notConnected  = errors.New("Not connected")
	notAssociated = errors.New("Not associated")
)

const SessionStep = 10

type PfcpSimCfg struct {
	interfaceName string
	upfN3         string
	serverAddr    string
	state         int
	sim           *sim.PFCPClient
}

func NewPfcpSimCfg(iface, upfN3, serverAddr string) *PfcpSimCfg {
	return &PfcpSimCfg{
		interfaceName: iface,
		upfN3:         upfN3,
		serverAddr:    serverAddr,
		state:         SVC_INIT,
		sim:           nil,
	}
}

func (c *PfcpSimCfg) InitPFCPSim() error {
	pfcpsim.NewPFCPSimService(c.interfaceName)
	pfcpsim.SetRemotePeer(c.serverAddr)
	pfcpsim.SetUpfN3(c.upfN3)

	var err error

	if err = pfcpsim.ConnectPFCPSim(); err != nil {
		return err
	} else {
		c.state = SVC_CONNECTED
		c.sim = pfcpsim.GetSimulator()
	}

	return nil
}

func (c *PfcpSimCfg) TerminatePFCPSim() error {
	return pfcpsim.DisconnectPFCPSim()
}

func (c *PfcpSimCfg) Associate() error {
	switch c.state {
	case SVC_INIT:
		return notConnected
	case SVC_CONNECTED:
		err := c.sim.SetupAssociation()
		if err != nil {
			return err
		}
	}

	c.state = SVC_ASSOCIATED

	return nil
}

func (c *PfcpSimCfg) Deassociate() error {
	switch c.state {
	case SVC_INIT:
		return notConnected
	case SVC_CONNECTED:
		return notAssociated
	case SVC_ASSOCIATED:
		err := c.sim.TeardownAssociation()
		if err != nil {
			return err
		}

		c.sim.DisconnectN4()
	}

	c.state = SVC_INIT

	return nil
}

type SessionIEs struct {
}

func (c *PfcpSimCfg) CreateSession(baseID int,
	pdrBuilder, qerBuilder, farBuilder, urrBuilder int,
	fuzz uint) error {
	var uePool = "10.60.0.0/15"

	appFilters := []string{""}
	count := 1

	lastUEAddr, _, err := net.ParseCIDR(uePool)
	if err != nil {
		errMsg := fmt.Sprintf(" Could not parse Address Pool: %v", err)
		log.Error(errMsg)

		return status.Error(codes.Aborted, errMsg)
	}

	var qfi uint8 = 6

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
				FuzzIE(qerBuilder, fuzz).
				Build(),
		}

		// create as many PDRs, FARs and App QERs as the number of app filters provided through pfcpctl
		ID := uint16(i)

		for _, appFilter := range appFilters {
			SDFFilter, gateStatus, precedence, err := pfcpsim.ParseAppFilter(appFilter)
			if err != nil {
				return status.Error(codes.Aborted, err.Error())
			}

			log.Infof("Successfully parsed application filter. SDF Filter: %v", SDFFilter)

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
				WithMeasurementPeriod(1*time.Second).
				WithReportingTrigger(session.ReportingTrigger{
					Flags: session.RPT_TRIG_PERIO,
				}).
				FuzzIE(urrBuilder, fuzz).
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
				FuzzIE(urrBuilder, fuzz).
				Build()

			urrs = append(urrs, urr)

			uplinkPDR := session.NewPDRBuilder().
				WithID(uplinkPdrID).
				WithMethod(session.Create).
				WithTEID(uplinkTEID).
				WithFARID(uplinkFarID).
				AddQERID(sessQerID).
				AddQERID(uplinkAppQerID).
				WithN3Address(c.upfN3).
				WithSDFFilter(SDFFilter).
				WithPrecedence(precedence).
				MarkAsUplink().
				FuzzIE(pdrBuilder, fuzz).
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
				FuzzIE(pdrBuilder, fuzz).
				BuildPDR()

			pdrs = append(pdrs, uplinkPDR)
			pdrs = append(pdrs, downlinkPDR)

			uplinkFAR := session.NewFARBuilder().
				WithID(uplinkFarID).
				WithAction(session.ActionForward).
				WithDstInterface(ieLib.DstInterfaceCore).
				WithMethod(session.Create).
				FuzzIE(farBuilder, fuzz).
				BuildFAR()

			downlinkFAR := session.NewFARBuilder().
				WithID(downlinkFarID).
				WithAction(session.ActionDrop).
				WithMethod(session.Create).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithZeroBasedOuterHeaderCreation().
				FuzzIE(farBuilder, fuzz).
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
				FuzzIE(qerBuilder, fuzz).
				Build()

			downlinkAppQER := session.NewQERBuilder().
				WithID(downlinkAppQerID).
				WithMethod(session.Create).
				WithQFI(qfi).
				WithUplinkMBR(50000).
				WithDownlinkMBR(30000).
				WithGateStatus(gateStatus).
				FuzzIE(qerBuilder, fuzz).
				Build()

			qers = append(qers, uplinkAppQER)
			qers = append(qers, downlinkAppQER)

			ID += 2
		}

		sess, err := c.sim.EstablishSession(pdrs, fars, qers, urrs)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		sim.InsertSession(i, sess)
	}

	infoMsg := fmt.Sprintf("%v sessions were established using %v as baseID ", count, baseID)
	log.Info(infoMsg)

	return nil
}

func (c *PfcpSimCfg) ModifySession(baseID int,
	farBuilder, urrBuilder int,
	fuzz uint) error {
	var actions uint8

	actions |= session.ActionNotify
	count := 1
	nodeBaddress := "192.168.0.1"
	appFilters := []string{""}

	for i := baseID; i < (count*SessionStep + baseID); i = i + SessionStep {
		var newFARs []*ieLib.IE

		var newURRs []*ieLib.IE

		ID := uint32(i + 1)
		teid := uint32(i + 1)

		for j := 0; j < len(appFilters); j++ {
			downlinkFAR := session.NewFARBuilder().
				WithID(ID). // Same FARID that was generated in create sessions
				WithMethod(session.Update).
				WithAction(actions).
				WithDstInterface(ieLib.DstInterfaceAccess).
				WithTEID(teid).
				WithDownlinkIP(nodeBaddress).
				FuzzIE(farBuilder, fuzz).
				BuildFAR()

			newFARs = append(newFARs, downlinkFAR)

			urrId := uint32(ID)
			urr := session.NewURRBuilder().
				WithID(urrId).
				WithMethod(session.Update).
				WithMeasurementPeriod(1*time.Second).
				FuzzIE(urrBuilder, fuzz).
				Build()

			newURRs = append(newURRs, urr)

			ID += 2
		}

		sess, ok := sim.GetSession(i)
		if !ok {
			errMsg := fmt.Sprintf("Could not retrieve session with index %v", i)
			log.Error(errMsg)

			return status.Error(codes.Internal, errMsg)
		}

		err := c.sim.ModifySession(sess, nil, newFARs, nil, newURRs)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	infoMsg := fmt.Sprintf("%v sessions were modified", count)
	log.Info(infoMsg)

	return nil
}

func (c *PfcpSimCfg) DeleteSession(baseID int) error {
	count := 1

	if sim.GetActiveSessionNum() < count {
		err := sim.NewNotEnoughSessionsError()
		log.Error(err)

		return status.Error(codes.Aborted, err.Error())
	}

	for i := baseID; i < (count*SessionStep + baseID); i = i + SessionStep {
		sess, ok := sim.GetSession(i)
		if !ok {
			errMsg := "Session was nil. Check baseID"
			log.Error(errMsg)

			return status.Error(codes.Aborted, errMsg)
		}

		err := c.sim.DeleteSession(sess)
		if err != nil {
			log.Error(err.Error())
			return status.Error(codes.Aborted, err.Error())
		}
		// remove from activeSessions
		sim.RemoveSession(i)
	}

	infoMsg := fmt.Sprintf("%v sessions deleted; activeSessions: %v", count, sim.GetActiveSessionNum())
	log.Info(infoMsg)

	return nil
}
