package pfcpsim

import "github.com/wmnsk/go-pfcp/ie"

type PFCPSession struct {
	LocalSEID uint64
	PeerSEID  uint64

	PDRs []*ie.IE
	FARs []*ie.IE

	QERs []*ie.IE
}

func NewSession() *PFCPSession {
	return &PFCPSession{
		LocalSEID: 0, // Updated by PFCPClient when establishing session
		PeerSEID:  0,  // Updated later by PFCPClient when received F-SEID IE from peer
		PDRs:      make([]*ie.IE, 0),
		FARs:      make([]*ie.IE, 0),
		QERs:      make([]*ie.IE, 0),
	}
}
