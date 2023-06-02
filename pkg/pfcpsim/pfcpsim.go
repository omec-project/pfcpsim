// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package pfcpsim

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/wmnsk/go-pfcp/ie"
	ieLib "github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"
)

const (
	PFCPStandardPort       = 8805
	DefaultHeartbeatPeriod = 5
	DefaultResponseTimeout = 5 * time.Second
)

var (
	activeSessions     = make(map[int]*PFCPSession, 0)
	lockActiveSessions = new(sync.Mutex)
)

func GetActiveSessionNum() int {
	lockActiveSessions.Lock()
	defer lockActiveSessions.Unlock()

	return len(activeSessions)
}

func InsertSession(index int, session *PFCPSession) {
	lockActiveSessions.Lock()
	defer lockActiveSessions.Unlock()

	activeSessions[index] = session
}

func GetSession(index int) (*PFCPSession, bool) {
	lockActiveSessions.Lock()
	defer lockActiveSessions.Unlock()
	element, ok := activeSessions[index]
	return element, ok
}

func GetSessionByLocalSEID(seid uint64) (*PFCPSession, bool) {
	lockActiveSessions.Lock()
	defer lockActiveSessions.Unlock()
	for _, session := range activeSessions {
		if session.localSEID == seid {
			return session, true
		}
	}
	return nil, false
}

func RemoveSession(index int) {
	lockActiveSessions.Lock()
	defer lockActiveSessions.Unlock()

	delete(activeSessions, index)
}

// PFCPClient enables to simulate a client sending PFCP messages towards the UPF.
// It provides two usage modes:
//   - 1st mode enables high-level PFCP operations (e.g., SetupAssociation())
//   - 2nd mode gives a user more control over PFCP sequence flow
//     and enables send and receive of individual messages (e.g., SendAssociationSetupRequest(), PeekNextResponse())
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

	localAddr  string
	remoteAddr string
	conn       *net.UDPConn

	// responseTimeout timeout to wait for PFCP response (default: 5 seconds)
	responseTimeout time.Duration
}

func NewPFCPClient(localAddr string) *PFCPClient {
	client := &PFCPClient{
		sequenceNumber:  0,
		localAddr:       localAddr,
		responseTimeout: DefaultResponseTimeout,
	}

	client.ctx = context.Background()
	client.heartbeatsChan = make(chan *message.HeartbeatResponse)
	client.recvChan = make(chan message.Message)

	return client
}

func (c *PFCPClient) SetPFCPResponseTimeout(timeout time.Duration) {
	c.responseTimeout = timeout
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

	raddr, err := net.ResolveUDPAddr("udp", c.remoteAddr)
	if err != nil {
		return err
	}

	if _, err := c.conn.WriteTo(b, raddr); err != nil {
		return err
	}

	return nil
}

func (c *PFCPClient) receiveFromN4(ctx context.Context) {
	buf := make([]byte, 1500)

	for {
		select {
		case <-ctx.Done():
			if c.cancelHeartbeats != nil {
				c.cancelHeartbeats()
			}
			c.conn.Close()
			return
		default:
			n, _, err := c.conn.ReadFrom(buf)
			if err != nil {
				continue
			}

			msg, err := message.Parse(buf[:n])
			if err != nil {
				continue
			}

			switch msg := msg.(type) {
			case *message.HeartbeatResponse:
				c.heartbeatsChan <- msg

			case *message.SessionReportRequest:
				if c.handleSessionReportRequest(msg) {
					continue
				}
			default:
				c.recvChan <- msg
			}
		}
	}
}

func (c *PFCPClient) ConnectN4(ctx context.Context, remoteAddr string) error {
	addr := fmt.Sprintf("%s:%d", remoteAddr, PFCPStandardPort)

	if host, port, err := net.SplitHostPort(remoteAddr); err == nil {
		// remoteAddr contains also a port. Use provided port instead of PFCPStandardPort
		addr = fmt.Sprintf("%s:%s", host, port)
	}

	c.remoteAddr = addr

	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.localAddr, PFCPStandardPort))
	if err != nil {
		return err
	}

	rxconn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return err
	}

	c.conn = rxconn

	go c.receiveFromN4(ctx)

	return nil
}

func (c *PFCPClient) DisconnectN4() {
	if c.cancelHeartbeats != nil {
		c.cancelHeartbeats()
	}

	c.conn.Close()
}

func (c *PFCPClient) PeekNextHeartbeatResponse() (*message.HeartbeatResponse, error) {
	select {
	case msg := <-c.heartbeatsChan:
		return msg, nil
	case <-time.After(c.responseTimeout):
		return nil, NewTimeoutExpiredError()
	}
}

// PeekNextResponse can be used to wait for a next PFCP message from a peer.
// It's a blocking operation, which is timed out after c.responseTimeout period (5 seconds by default).
// Use SetPFCPResponseTimeout() to configure a custom timeout.
func (c *PFCPClient) PeekNextResponse() (message.Message, error) {
	var resMsg message.Message

	delay := time.NewTimer(c.responseTimeout)

	for {
		select {
		case resMsg := <-c.recvChan:
			if !delay.Stop() {
				<-delay.C
			}
			return resMsg, nil
		case <-delay.C:
			return resMsg, NewTimeoutExpiredError()
		}
	}
}

// MsgTypeSessionReportRequest: sent by the UP function to the CP function to report information related to an PFCP session
// MsgTypeSessionReportResponse: sent by the CP function to the UP function as a reply to the Session Report Request.
func (c *PFCPClient) handleSessionReportRequest(msg *message.SessionReportRequest) bool {
	if msg.MessageType() == message.MsgTypeSessionReportRequest {
		fmt.Println("Session Report Request received")
		err := c.sendSessionReportResponse(msg.Sequence(),
			msg.Header.SEID)
		if err != nil {
			fmt.Println("Error sending Session Report Response")
		}
		return true
	}
	return false
}

func (c *PFCPClient) sendSessionReportResponse(seq uint32, seid uint64) error {
	var rseid uint64
	sess, ok := GetSessionByLocalSEID(seid)
	if !ok {
		rseid = 0
		fmt.Println("Session not found")
	} else {
		rseid = sess.peerSEID
	}
	res := message.NewSessionReportResponse(0, 0, rseid, seq, 0,
		ie.NewCause(ie.CauseRequestAccepted))

	return c.sendMsg(res)
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

// SendAssociationTeardownRequest sends PFCP Teardown Request towards a peer.
// A caller should make sure that the PFCP connection is established before invoking this function.
func (c *PFCPClient) SendAssociationTeardownRequest(ie ...*ieLib.IE) error {
	teardownReq := message.NewAssociationReleaseRequest(0,
		ieLib.NewNodeID(c.conn.RemoteAddr().String(), "", ""),
	)

	teardownReq.IEs = append(teardownReq.IEs, ie...)

	return c.sendMsg(teardownReq)
}

func (c *PFCPClient) SendHeartbeatRequest() error {
	hbReq := message.NewHeartbeatRequest(
		c.getNextSequenceNumber(),
		ieLib.NewRecoveryTimeStamp(time.Now()),
		ieLib.NewSourceIPAddress(net.ParseIP(c.localAddr), nil, 0),
	)

	return c.sendMsg(hbReq)
}

func (c *PFCPClient) SendSessionEstablishmentRequest(pdrs []*ieLib.IE, fars []*ieLib.IE,
	qers []*ieLib.IE, urrs []*ieLib.IE) error {
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
	estReq.CreateURR = append(estReq.CreateURR, urrs...)

	return c.sendMsg(estReq)
}

func (c *PFCPClient) SendSessionModificationRequest(PeerSEID uint64, pdrs []*ieLib.IE, qers []*ieLib.IE, fars []*ieLib.IE, urrs []*ieLib.IE) error {
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
	modifyReq.UpdateURR = append(modifyReq.UpdateURR, urrs...)

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

	_, err = c.PeekNextHeartbeatResponse()
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

	resp, err := c.PeekNextResponse()
	if err != nil {
		return err
	}

	assocResp, ok := resp.(*message.AssociationSetupResponse)
	if !ok {
		return NewInvalidResponseError(fmt.Errorf("unexpected response type"))
	}

	cause, err := assocResp.Cause.Cause()
	if err != nil {
		return NewInvalidResponseError(err)
	}

	if cause != ieLib.CauseRequestAccepted {
		return NewInvalidResponseError(fmt.Errorf("association setup failed with cause %d", cause))
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

	err := c.SendAssociationTeardownRequest()
	if err != nil {
		return err
	}

	resp, err := c.PeekNextResponse()
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
func (c *PFCPClient) EstablishSession(pdrs []*ieLib.IE, fars []*ieLib.IE,
	qers []*ieLib.IE, urrs []*ieLib.IE) (*PFCPSession, error) {
	if !c.isAssociationActive {
		return nil, NewAssociationInactiveError()
	}

	err := c.SendSessionEstablishmentRequest(pdrs, fars, qers, urrs)
	if err != nil {
		return nil, err
	}

	resp, err := c.PeekNextResponse()
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

func (c *PFCPClient) ModifySession(sess *PFCPSession, pdrs []*ieLib.IE, fars []*ieLib.IE,
	qers []*ieLib.IE, urrs []*ieLib.IE) error {
	if !c.isAssociationActive {
		return NewAssociationInactiveError()
	}

	err := c.SendSessionModificationRequest(sess.peerSEID, pdrs, fars, qers, urrs)
	if err != nil {
		return err
	}

	resp, err := c.PeekNextResponse()
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

	resp, err := c.PeekNextResponse()
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
