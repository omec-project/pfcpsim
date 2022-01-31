package session

import (
	"net"

	"github.com/wmnsk/go-pfcp/ie"
)

type pdrBuilder struct {
	precedence uint32
	method     IEMethod
	sdfFilter  string
	id         uint16
	teid       uint32

	qerIDs []*ie.IE
	farIDs []*ie.IE

	ueAddress string
	n3Address string

	direction direction
}

func NewPDRBuilder() *pdrBuilder {
	return &pdrBuilder{
		qerIDs: make([]*ie.IE, 0),
		farIDs: make([]*ie.IE, 0),
	}
}

func (b *pdrBuilder) WithPrecedence(precedence uint32) *pdrBuilder {
	b.precedence = precedence
	return b
}

func (b *pdrBuilder) WithSDFFilter(filter string) *pdrBuilder {
	b.sdfFilter = filter
	return b
}

func (b *pdrBuilder) WithID(id uint16) *pdrBuilder {
	b.id = id
	return b
}

func (b *pdrBuilder) WithTEID(teid uint32) *pdrBuilder {
	b.teid = teid
	return b
}

func (b *pdrBuilder) WithMethod(method IEMethod) *pdrBuilder {
	b.method = method
	return b
}

func (b *pdrBuilder) WithN3Address(n3Address string) *pdrBuilder {
	b.n3Address = n3Address
	return b
}

func (b *pdrBuilder) WithUEAddress(ueAddress string) *pdrBuilder {
	b.ueAddress = ueAddress
	return b
}

func (b *pdrBuilder) AddQERID(qerID uint32) *pdrBuilder {
	b.qerIDs = append(b.qerIDs, ie.NewQERID(qerID))
	return b
}

func (b *pdrBuilder) AddFARID(farID uint32) *pdrBuilder {
	b.farIDs = append(b.farIDs, ie.NewFARID(farID))
	return b
}

func (b *pdrBuilder) MarkAsDownlink() *pdrBuilder {
	b.direction = downlink
	return b
}

func (b *pdrBuilder) MarkAsUplink() *pdrBuilder {
	b.direction = uplink
	return b
}

func (b *pdrBuilder) validate() {
	if b.direction == notSet {
		panic("Tried building a PDR without marking it as uplink or downlink")
	}

	if len(b.qerIDs) == 0 {
		panic("Tried building PDR without providing QER IDs")
	}

	if len(b.farIDs) == 0 {
		panic("Tried building PDR without providing FAR IDs")
	}

	if b.direction == downlink {
		if b.ueAddress == "" {
			panic("Tried building downlink PDR without setting the UE IP address")
		}
	}

	if b.direction == uplink {
		if b.n3Address == "" {
			panic("Tried building uplink PDR without setting the N3Address")
		}
		if b.teid == 0 {
			panic("Tried building uplink PDR without setting the TEID")
		}
	}
}

// BuildPDR returns by default an UplinkFAR.
// Returns a DownlinkFAR if MarkAsDownlink was invoked.
func (b *pdrBuilder) BuildPDR() *ie.IE {
	b.validate()

	createFunc := ie.NewCreatePDR
	if b.method == Update {
		createFunc = ie.NewUpdatePDR
	}

	if b.direction == downlink {
		pdr := createFunc(
			ie.NewPDRID(b.id),
			ie.NewPrecedence(b.precedence),
			ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceCore),
				ie.NewUEIPAddress(0x2, b.ueAddress, "", 0, 0),
				ie.NewSDFFilter(b.sdfFilter, "", "", "", 1),
			),
		)

		pdr.Add(b.farIDs...)
		pdr.Add(b.qerIDs...)
		return pdr
	}

	// UplinkPDR
	pdr := createFunc(
		ie.NewPDRID(b.id),
		ie.NewPrecedence(b.precedence),
		ie.NewPDI(
			ie.NewSourceInterface(ie.SrcInterfaceAccess),
			ie.NewFTEID(0x01, b.teid, net.ParseIP(b.n3Address), nil, 0),
			ie.NewSDFFilter(b.sdfFilter, "", "", "", 1),
		),
		ie.NewOuterHeaderRemoval(0, 0),
	)

	pdr.Add(b.farIDs...)
	pdr.Add(b.qerIDs...)

	return pdr

}
