// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package pfcpsim

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

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
	// keeps the current number of active PFCP sessions
	// it is also used as F-SEID
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

// ListenN4 allows for UP initiated PFCP associations.
func (c *PFCPClient) ListenN4() error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: PFCPStandardPort})
	if err != nil {
		return err
	}

	c.conn = conn

	go c.receiveFromN4()

	return nil
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
		return nil, NewTimeoutExpiredError()
	}
}

func (c *PFCPClient) PeekNextResponse(timeout time.Duration) (message.Message, error) {
	select {
	case msg := <-c.recvChan:
		return msg, nil
	case <-time.After(timeout * time.Second):
		return nil, NewTimeoutExpiredError()
	}
}

func (c *PFCPClient) SendAssociationSetupRequest(ie ...*ieLib.IE) error {
	c.resetSequenceNumber()

	assocReq := message.NewAssociationSetupRequest(
		c.getNextSequenceNumber(),
		ieLib.NewRecoveryTimeStamp(time.Now()),
		ieLib.NewNodeID(c.localAddr, "", ""),
	)

	assocReq.IEs = append(assocReq.IEs, ie...)

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

func (c *PFCPClient) SendSessionEstablishmentRequest(pdrs []*ieLib.IE, fars []*ieLib.IE, qers []*ieLib.IE) error {
	estReq := message.NewSessionEstablishmentRequest(
		0,
		0,
		0,
		c.getNextSequenceNumber(),
		0,
		ieLib.NewNodeID(c.localAddr, "", ""),
		ieLib.NewFSEID(c.getNextFSEID(), net.ParseIP(c.localAddr), nil),
		ieLib.NewPDNType(ieLib.PDNTypeIPv4),
	)
	estReq.CreatePDR = append(estReq.CreatePDR, pdrs...)
	estReq.CreateFAR = append(estReq.CreateFAR, fars...)
	estReq.CreateQER = append(estReq.CreateQER, qers...)

	return c.sendMsg(estReq)
}

func (c *PFCPClient) SendSessionModificationRequest(PeerSEID uint64, pdrs []*ieLib.IE, qers []*ieLib.IE, fars []*ieLib.IE) error {
	modifyReq := message.NewSessionModificationRequest(
		0,
		0,
		PeerSEID,
		c.getNextSequenceNumber(),
		0,
	)

	modifyReq.UpdatePDR = append(modifyReq.UpdatePDR, pdrs...)
	modifyReq.UpdateFAR = append(modifyReq.UpdateFAR, fars...)
	modifyReq.UpdateQER = append(modifyReq.UpdateQER, qers...)

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
		return NewInvalidResponseError()
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
		return NewAssociationInactiveError()
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
		return NewInvalidResponseError()
	}

	if c.cancelHeartbeats != nil {
		c.cancelHeartbeats()
	}

	c.setAssociationStatus(false)

	return nil
}

// EstablishSession sends PFCP Session Establishment Request and waits for PFCP Session Establishment Response.
// Returns a pointer to a new PFCPSession. Returns error if the process fails at any stage.
func (c *PFCPClient) EstablishSession(pdrs []*ieLib.IE, fars []*ieLib.IE, qers []*ieLib.IE) (*PFCPSession, error) {
	if !c.isAssociationActive {
		return nil, NewAssociationInactiveError()
	}

	err := c.SendSessionEstablishmentRequest(pdrs, fars, qers)
	if err != nil {
		return nil, err
	}

	resp, err := c.PeekNextResponse(5)
	if err != nil {
		return nil, NewTimeoutExpiredError(err)
	}

	estResp, ok := resp.(*message.SessionEstablishmentResponse)
	if !ok {
		return nil, NewInvalidResponseError(err)
	}

	if cause, err := estResp.Cause.Cause(); err != nil || cause != ieLib.CauseRequestAccepted {
		return nil, NewInvalidCauseError(err)
	}

	remoteSEID, err := estResp.UPFSEID.FSEID()
	if err != nil {
		return nil, err
	}

	sess := &PFCPSession{
		localSEID: c.lastFSEID,
		peerSEID:  remoteSEID.SEID,
	}

	return sess, nil
}

func (c *PFCPClient) ModifySession(sess *PFCPSession, pdrs []*ieLib.IE, fars []*ieLib.IE, qers []*ieLib.IE) error {
	if !c.isAssociationActive {
		return NewAssociationInactiveError()
	}

	err := c.SendSessionModificationRequest(sess.peerSEID, pdrs, fars, qers)
	if err != nil {
		return err
	}

	resp, err := c.PeekNextResponse(5)
	if err != nil {
		return NewTimeoutExpiredError(err)
	}

	modRes, ok := resp.(*message.SessionModificationResponse)
	if !ok {
		return NewInvalidResponseError(err)
	}

	if cause, err := modRes.Cause.Cause(); err != nil || cause != ieLib.CauseRequestAccepted {
		return NewInvalidCauseError(err)
	}

	return nil
}

// DeleteSession sends Session Deletion Request for each session and awaits for PFCP Session Deletion Response.
// Returns error if the process fails at any stage.
func (c *PFCPClient) DeleteSession(sess *PFCPSession) error {
	err := c.SendSessionDeletionRequest(sess.localSEID, sess.peerSEID)
	if err != nil {
		return err
	}

	resp, err := c.PeekNextResponse(5)
	if err != nil {
		return err
	}

	delResp, ok := resp.(*message.SessionDeletionResponse)
	if !ok {
		return NewInvalidResponseError()
	}

	if cause, err := delResp.Cause.Cause(); err != nil || cause != ieLib.CauseRequestAccepted {
		return NewInvalidCauseError(err)
	}

	return nil
}
