package session

import (
	"github.com/wmnsk/go-pfcp/ie"
)

type Session struct {
	LocalSEID uint64

	PeerSEID uint64

	UplinkPDRs   []*ie.IE
	DownlinkPDRs []*ie.IE

	UplinkFARs   []*ie.IE
	DownlinkFARs []*ie.IE

	QERs []*ie.IE
}

func NewSession() *Session {
	return &Session{
		LocalSEID:    0, // Updated by PFCPClient when establishing session
		PeerSEID:     0, // Updated later by PFCPClient when received F-SEID IE from peer
		UplinkPDRs:   make([]*ie.IE, 0),
		DownlinkPDRs: make([]*ie.IE, 0),
		UplinkFARs:   make([]*ie.IE, 0),
		DownlinkFARs: make([]*ie.IE, 0),
		QERs:         make([]*ie.IE, 0),
	}
}

func (s *Session) ClearSentRules() {
	s.UplinkPDRs = make([]*ie.IE, 0)
	s.DownlinkPDRs = make([]*ie.IE, 0)
	s.UplinkFARs = make([]*ie.IE, 0)
	s.DownlinkFARs = make([]*ie.IE, 0)
	s.QERs = make([]*ie.IE, 0)
}

func (s *Session) IsActive() bool {
	return s.PeerSEID != 0
}
