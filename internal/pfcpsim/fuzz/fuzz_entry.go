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
	if err := pfcpsim.ConnectPFCPSim(); err != nil {
		return fmt.Errorf("Could not connect to remote peer :%v", err)
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
		return errors.New("Not connected")
	case SVC_CONNECTED:
		err := c.sim.SetupAssociation()
		if err != nil {
			return err
		}
		c.state = SVC_ASSOCIATED
	case SVC_ASSOCIATED:
		return fmt.Errorf("Already associated")
	}
	return nil
}

func (c *PfcpSimCfg) Deassociate() error {
	switch c.state {
	case SVC_INIT:
		return fmt.Errorf("Not connected")
	case SVC_CONNECTED:
		return fmt.Errorf("Not associated")
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

type CreateSessionIEs struct{}

func (c *PfcpSimCfg) CreateSession() error {
	var AppFilters []string
	uePool := "10.60.0.0/15"

	baseID := 2
	count := 1

	if len(AppFilters) == 0 {
		AppFilters = append(AppFilters, "")
	}
	lastUEAddr, _, err := net.ParseCIDR(uePool)
	if err != nil {
		errMsg := fmt.Sprintf(" Could not parse Address Pool: %v", err)
		log.Error(errMsg)
		return status.Error(codes.Aborted, errMsg)
	}

	var qfi uint8 = 6
	SessionStep := 10

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

		for _, appFilter := range AppFilters {
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
				WithN3Address(c.upfN3).
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
