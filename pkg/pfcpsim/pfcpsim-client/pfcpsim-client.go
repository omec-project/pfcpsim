// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package pfcpsim_client

import (
	"github.com/c-robinson/iplib"
	"github.com/omec-project/pfcpsim/pkg/pfcpsim"
	log "github.com/sirupsen/logrus"
	"github.com/wmnsk/go-pfcp/ie"
	"net"
)

const (
	ActionForward uint8 = 0x2
	ActionDrop    uint8 = 0x1
	ActionBuffer  uint8 = 0x4
	ActionNotify  uint8 = 0x8
)

type PFCPSimClient struct {
	ueAddressPool string
	nodeBAddress  string
	upfAddress    string

	activeSessions int

	lastUEAddress net.IP

	client *pfcpsim.PFCPClient
}

func NewPFCPSimClient(lAddr string, ueAddressPool string, nodeBAddress string, upfAddress string) *PFCPSimClient {

	pfcpClient := pfcpsim.NewPFCPClient(lAddr)

	return &PFCPSimClient{
		ueAddressPool: ueAddressPool,
		nodeBAddress:  nodeBAddress,
		upfAddress:    upfAddress,
		client:        pfcpClient,
	}
}

func (c *PFCPSimClient) Disconnect() {
	c.client.DisconnectN4()
	log.Infof("PFCP client Disconnected")
}

func (c *PFCPSimClient) Connect(remoteAddress string) error {
	err := c.client.ConnectN4(remoteAddress)
	if err != nil {
		return err
	}

	log.Infof("PFCP client is connected")
	return nil
}

// TeardownAssociation uses the PFCP client to tearing down an already established association.
// If called while no association is established by PFCP client, the latter will return an error
func (c *PFCPSimClient) TeardownAssociation() {
	err := c.client.TeardownAssociation()
	if err != nil {
		log.Errorf("Error while tearing down association: %v", err)
		return
	}

	log.Infoln("Teardown association completed")
}

// SetupAssociation uses the PFCP client to establish an association,
// logging its success by checking PFCPclient.IsAssociationAlive
func (c *PFCPSimClient) SetupAssociation() {
	err := c.client.SetupAssociation()
	if err != nil {
		log.Errorf("Error while setting up association: %v", err)
		return
	}

	if !c.client.IsAssociationAlive() {
		log.Errorf("Error while peeking heartbeat response: %v", err)
		return
	}

	log.Infof("Setup association completed")
}

// getNextUEAddress retrieves the next available IP address from ueAddressPool
func (c *PFCPSimClient) getNextUEAddress() net.IP {
	// TODO handle case net address is full
	if c.lastUEAddress == nil {
		ueIpFromPool, _, _ := net.ParseCIDR(c.ueAddressPool)
		c.lastUEAddress = iplib.NextIP(ueIpFromPool)

		return c.lastUEAddress

	} else {
		c.lastUEAddress = iplib.NextIP(c.lastUEAddress)
		return c.lastUEAddress
	}
}

// InitializeSessions create 'count' sessions incrementally.
// Once created, the sessions are established through PFCP client.
func (c *PFCPSimClient) InitializeSessions(count int) {
	baseID := c.activeSessions + 1

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

		uplinkAppQerID := uint32(i)
		downlinkAppQerID := uint32(i + 1)

		pdrs := []*ie.IE{
			pfcpsim.NewUplinkPDR(pfcpsim.Create, uplinkPdrID, uplinkTEID, c.upfAddress, uplinkFarID, sessQerID, uplinkAppQerID),
			pfcpsim.NewDownlinkPDR(pfcpsim.Create, dowlinkPdrID, c.getNextUEAddress().String(), downlinkFarID, sessQerID, downlinkAppQerID),
		}

		fars := []*ie.IE{
			pfcpsim.NewUplinkFAR(pfcpsim.Create, uplinkFarID, ActionForward),
			pfcpsim.NewDownlinkFAR(pfcpsim.Create, downlinkFarID, ActionDrop, downlinkTEID, c.nodeBAddress),
		}

		qers := []*ie.IE{
			// session QER
			pfcpsim.NewQER(pfcpsim.Create, sessQerID, 0x09, 500000, 500000, 0, 0),
			// application QER
			pfcpsim.NewQER(pfcpsim.Create, appQerID, 0x08, 50000, 50000, 30000, 30000),
		}

		err := c.client.EstablishSession(pdrs, fars, qers)
		if err != nil {
			log.Errorf("Error while establishing sessions: %v", err)
			return
		}

		// TODO show session's F-SEID
		c.activeSessions++
		log.Infof("Created session")
	}

}

// DeleteAllSessions uses the PFCP client DeleteAllSessions. If failure happens at any stage,
// an error is logged.
func (c *PFCPSimClient) DeleteAllSessions() {
	err := c.client.DeleteAllSessions()
	if err != nil {
		log.Errorf("Error while deleting sessions: %v", err)
		return
	}

	log.Infof("Deleted all sessions")
}
