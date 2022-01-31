// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package pfcpsim

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/omec-project/pfcpsim/pkg/pfcpsim/session"
	ieLib "github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"
)

const (
	PFCPStandardPort       = 8805
	DefaultHeartbeatPeriod = 5
)

// PFCPClient enables to simulate a client sending PFCP messages towards the UPF.
// It provides two usage modes:
// - 1st mode enables high-level PFCP operations (e.g., SetupAssociation())
// - 2nd mode gives a user more control over PFCP sequence flow
//   and enables send and receive of individual messages (e.g., SendAssociationSetupRequest(), PeekNextResponse())
type PFCPClient struct {
	// keeps track of active PFCP sessions
	activeSessions map[uint64]*session.Session
	// keeps track of last created FSEID
	lastFSEID uint64

	aliveLock           sync.Mutex
	isAssociationActive bool

	ctx              context.Context
	cancelHeartbeats context.CancelFunc

	heartbeatsChan chan *message.HeartbeatResponse
	recvChan       chan message.Message

	sequenceNumber uint32
	seqNumLock     sync.Mutex

	localAddr string
	conn      *net.UDPConn
}

func NewPFCPClient(localAddr string) *PFCPClient {
	client := &PFCPClient{
		sequenceNumber: 0,
		localAddr:      localAddr,
		activeSessions: make(map[uint64]*session.Session),
	}
	client.ctx = context.Background()
	client.heartbeatsChan = make(chan *message.HeartbeatResponse)
	client.recvChan = make(chan message.Message)
	return client
}

func (c *PFCPClient) getNextSequenceNumber() uint32 {
	c.seqNumLock.Lock()
	defer c.seqNumLock.Unlock()

	c.sequenceNumber++

	return c.sequenceNumber
}

func (c *PFCPClient) getNextFSEID() uint64 {
	c.lastFSEID++
	return c.lastFSEID
}

func (c *PFCPClient) resetSequenceNumber() {
	c.seqNumLock.Lock()
	defer c.seqNumLock.Unlock()

	c.sequenceNumber = 0
}

func (c *PFCPClient) setAssociationStatus(status bool) {
	c.aliveLock.Lock()
	defer c.aliveLock.Unlock()

	c.isAssociationActive = status
}

func (c *PFCPClient) sendMsg(msg message.Message) error {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return err
	}

	if _, err := c.conn.Write(b); err != nil {
		return err
	}

	return nil
}

func (c *PFCPClient) receiveFromN4() {
	buf := make([]byte, 1500)
	for {
		n, _, err := c.conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		msg, err := message.Parse(buf[:n])
		if err != nil {
			continue
		}

		if hbResp, ok := msg.(*message.HeartbeatResponse); ok {
			c.heartbeatsChan <- hbResp
		} else {
			c.recvChan <- msg
		}
	}
}

func (c *PFCPClient) ConnectN4(remoteAddr string) error {
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", remoteAddr, PFCPStandardPort))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	c.conn = conn

	go c.receiveFromN4()

	return nil
}

func (c *PFCPClient) DisconnectN4() {
	if c.cancelHeartbeats != nil {
		c.cancelHeartbeats()
	}
	c.conn.Close()
}

func (c *PFCPClient) PeekNextHeartbeatResponse(timeout time.Duration) (*message.HeartbeatResponse, error) {
	select {
	case msg := <-c.heartbeatsChan:
		return msg, nil
	case <-time.After(timeout * time.Second):
		return nil, errors.New("timeout waiting for response")
	}
}

func (c *PFCPClient) PeekNextResponse(timeout time.Duration) (message.Message, error) {
	select {
	case msg := <-c.recvChan:
		return msg, nil
	case <-time.After(timeout * time.Second):
		return nil, errors.New("timeout waiting for response")
	}
}

func (c *PFCPClient) SendAssociationSetupRequest(ie ...*ieLib.IE) error {
	c.resetSequenceNumber()

	assocReq := message.NewAssociationSetupRequest(
		c.getNextSequenceNumber(),
		ieLib.NewRecoveryTimeStamp(time.Now()),
		ieLib.NewNodeID(c.localAddr, "", ""),
	)

	for _, value := range ie {
		assocReq.IEs = append(assocReq.IEs, value)
	}

	return c.sendMsg(assocReq)
}

func (c *PFCPClient) SendHeartbeatRequest() error {
	hbReq := message.NewHeartbeatRequest(
		c.getNextSequenceNumber(),
		ieLib.NewRecoveryTimeStamp(time.Now()),
		ieLib.NewSourceIPAddress(net.ParseIP(c.localAddr), nil, 0),
	)

	return c.sendMsg(hbReq)
}

func (c *PFCPClient) SendSessionEstablishmentRequest(sessionToEstablish *session.Session) error {
	estReq := message.NewSessionEstablishmentRequest(
		0,
		0,
		0,
		c.getNextSequenceNumber(),
		0,
		ieLib.NewNodeID(c.localAddr, "", ""),
		ieLib.NewFSEID(sessionToEstablish.LocalSEID, net.ParseIP(c.localAddr), nil),
		ieLib.NewPDNType(ieLib.PDNTypeIPv4),
	)

	estReq.CreatePDR = append(estReq.CreatePDR, sessionToEstablish.UplinkPDRs...)
	estReq.CreatePDR = append(estReq.CreatePDR, sessionToEstablish.DownlinkPDRs...)

	estReq.CreateFAR = append(estReq.CreateFAR, sessionToEstablish.UplinkFARs...)
	estReq.CreateFAR = append(estReq.CreateFAR, sessionToEstablish.DownlinkFARs...)

	estReq.CreateQER = append(estReq.CreateQER, sessionToEstablish.QERs...)

	return c.sendMsg(estReq)
}

func (c *PFCPClient) SendSessionModificationRequest(far *ieLib.IE, PeerSEID uint64) error {
	// TODO in 5G mode also update PDR shall be sent
	modifyReq := message.NewSessionModificationRequest(
		0,
		0,
		PeerSEID,
		c.getNextSequenceNumber(),
		0,
		far,
	)

	return c.sendMsg(modifyReq)
}

func (c *PFCPClient) SendSessionDeletionRequest(localSEID uint64, remoteSEID uint64) error {
	delReq := message.NewSessionDeletionRequest(
		0,
		0,
		remoteSEID,
		c.getNextSequenceNumber(),
		0,
		ieLib.NewFSEID(localSEID, net.ParseIP(c.localAddr), nil),
	)

	return c.sendMsg(delReq)
}

func (c *PFCPClient) StartHeartbeats(stopCtx context.Context) {
	ticker := time.NewTicker(DefaultHeartbeatPeriod * time.Second)
	for {
		select {
		case <-stopCtx.Done():
			return
		case <-ticker.C:
			err := c.SendAndRecvHeartbeat()
			if err != nil {
				return
			}
		}
	}
}

func (c *PFCPClient) SendAndRecvHeartbeat() error {
	err := c.SendHeartbeatRequest()
	if err != nil {
		return err
	}

	_, err = c.PeekNextHeartbeatResponse(5)
	if err != nil {
		c.setAssociationStatus(false)
		return err
	}

	c.setAssociationStatus(true)

	return nil
}

// SetupAssociation sends PFCP Association Setup Request and waits for PFCP Association Setup Response.
// Returns error if the process fails at any stage.
func (c *PFCPClient) SetupAssociation() error {
	err := c.SendAssociationSetupRequest()
	if err != nil {
		return err
	}

	resp, err := c.PeekNextResponse(DefaultHeartbeatPeriod)
	if err != nil {
		return err
	}

	if _, ok := resp.(*message.AssociationSetupResponse); !ok {
		return fmt.Errorf("invalid message received, expected association setup response")
	}

	ctx, cancelFunc := context.WithCancel(c.ctx)
	c.cancelHeartbeats = cancelFunc

	c.setAssociationStatus(true)

	go c.StartHeartbeats(ctx)

	return nil
}

func (c *PFCPClient) IsAssociationAlive() bool {
	c.aliveLock.Lock()
	defer c.aliveLock.Unlock()

	return c.isAssociationActive
}

// TeardownAssociation tears down an already established association.
// If called while no association is established, an error is returned
func (c *PFCPClient) TeardownAssociation() error {
	if !c.IsAssociationAlive() {
		return errors.New("association does not exist")
	}

	ie1 := ieLib.NewNodeID(c.conn.RemoteAddr().String(), "", "")

	c.resetSequenceNumber()
	msg := message.NewAssociationReleaseRequest(0, ie1)

	err := c.sendMsg(msg)
	if err != nil {
		return err
	}

	resp, err := c.PeekNextResponse(5)
	if err != nil {
		return err
	}

	if _, ok := resp.(*message.AssociationReleaseResponse); !ok {
		return errors.New(fmt.Sprintf("received unexpected message: %v", resp.MessageTypeName()))
	}

	if c.cancelHeartbeats != nil {
		c.cancelHeartbeats()
	}
	c.setAssociationStatus(false)

	return nil
}

// EstablishSession sends PFCP Session Establishment Request and waits for PFCP Session Establishment Response.
// Returns error if the process fails at any stage.
func (c *PFCPClient) EstablishSession(s *session.Session) error {
	if !c.isAssociationActive {
		return fmt.Errorf("PFCP association is not active")
	}

	s.LocalSEID = c.getNextFSEID()

	err := c.SendSessionEstablishmentRequest(s)
	if err != nil {
		return err
	}

	resp, err := c.PeekNextResponse(5)
	if err != nil {
		return err
	}

	estResp, ok := resp.(*message.SessionEstablishmentResponse)
	if !ok {
		return fmt.Errorf("invalid message received, expected session establishment response")
	}

	if cause, err := estResp.Cause.Cause(); err != nil || cause != ieLib.CauseRequestAccepted {
		return fmt.Errorf("session establishment response returns invalid cause: %v", cause)
	}

	remoteSEID, err := estResp.UPFSEID.FSEID()
	if err != nil {
		return err
	}

	s.PeerSEID = remoteSEID.SEID

	c.activeSessions[c.lastFSEID] = s

	return nil
}

// ModifySessions sends a PFCP Session Modification Request for each active session,
// updating each downlinkFAR with action session.ActionDrop to session.ActionForward.
func (c *PFCPClient) ModifySessions(notifyCPFlag bool, bufferFlag bool, nodeBAddress *net.IP) error {
	counter := 0
	for _, activeSession := range c.activeSessions {
		for _, far := range activeSession.DownlinkFARs {
			action, err := far.ApplyAction()
			if err != nil {
				return err
			}

			if action != session.ActionDrop {
				// TODO is it ok to look only at action?
				// update only FARs with drop action
				continue
			}

			// recover FARID for update FAR
			farid, err := far.FARID()
			if err != nil {
				return err
			}

			// TODO handle notifyCP and buffer flag
			updateDownlinkFAR := session.NewFARBuilder().
				WithID(farid).
				WithMethod(session.Update).
				WithAction(session.ActionForward).
				WithTEID(activeSession.DownlinkTEID).
				WithDownlinkIP(nodeBAddress.String()).
				MarkAsDownlink().
				BuildFAR()

			err = c.SendSessionModificationRequest(updateDownlinkFAR, activeSession.PeerSEID)
			if err != nil {
				return err
			}

			modRes, err := c.PeekNextResponse(5)
			if err != nil {
				return err
			}

			modRes, ok := modRes.(*message.SessionModificationResponse)
			if !ok {
				return fmt.Errorf("invalid message received, expected session modification response")
			}

			counter++

			// Session was correctly modified -> save sent FAR in session's object
			// FIXME is it ok to clear sent Downlink FARs to avoid sending updateFARs to already updated flows?
			activeSession.DownlinkFARs = make([]*ieLib.IE, 2)

			activeSession.DownlinkFARs = append(activeSession.DownlinkFARs, updateDownlinkFAR)
		}
	}

	fmt.Printf("Sessions updated: %v\n", counter) // DEBUG remove

	return nil
}

// GetNumActiveSessions returns the number of active sessions.
func (c *PFCPClient) GetNumActiveSessions() int {
	return len(c.activeSessions)
}

// DeleteAllSessions sends Session Deletion Request for each session and awaits for PFCP Session Deletion Response.
// Returns error if the process fails at any stage.
func (c *PFCPClient) DeleteAllSessions() error {
	for _, activeSession := range c.activeSessions {

		err := c.SendSessionDeletionRequest(activeSession.LocalSEID, activeSession.PeerSEID)
		if err != nil {
			return err
		}

		resp, err := c.PeekNextResponse(5)
		if err != nil {
			return err
		}

		delResp, ok := resp.(*message.SessionDeletionResponse)
		if !ok {
			return fmt.Errorf("invalid message received, expected session deletion response")
		}

		if cause, err := delResp.Cause.Cause(); err != nil || cause != ieLib.CauseRequestAccepted {
			return fmt.Errorf("session deletion response returns invalid cause: %v", cause)
		}
		// remove session from active ones
		delete(c.activeSessions, activeSession.LocalSEID)
	}

	return nil
}
