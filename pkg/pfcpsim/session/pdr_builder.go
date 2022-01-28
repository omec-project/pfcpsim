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
	farID      uint32
	sessQerID  uint32
	appQerID   uint32

	ueAddress string
	n3Address string

	isDownlinkPDR bool
}

func NewPDRBuilder() *pdrBuilder {
	return &pdrBuilder{}
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

func (b *pdrBuilder) WithMEthod(method IEMethod) *pdrBuilder {
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

func (b *pdrBuilder) WithRulesIDs(farID uint32, sessionQERID uint32, appQERID uint32) *pdrBuilder {
	b.farID = farID
	b.sessQerID = sessionQERID
	b.appQerID = appQERID
	return b
}

func (b *pdrBuilder) MarkAsDownlink() *pdrBuilder {
	b.isDownlinkPDR = true
	return b
}

// BuildPDR returns by default an UplinkFAR.
// Returns a DownlinkFAR if MarkAsDownlink was invoked.
func (b *pdrBuilder) BuildPDR() *ie.IE {
	createFunc := ie.NewCreatePDR
	if b.method == Update {
		createFunc = ie.NewUpdatePDR
	}

	if b.isDownlinkPDR {
		return createFunc(
			ie.NewPDRID(b.id),
			ie.NewPrecedence(b.precedence),
			ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceCore),
				ie.NewUEIPAddress(0x2, b.ueAddress, "", 0, 0),
				ie.NewSDFFilter(b.sdfFilter, "", "", "", 1),
			),
			ie.NewFARID(b.farID),
			ie.NewQERID(b.sessQerID),
			ie.NewQERID(b.appQerID),
		)
	}

	// UplinkPDR
	return createFunc(
		ie.NewPDRID(b.id),
		ie.NewPrecedence(b.precedence),
		ie.NewPDI(
			ie.NewSourceInterface(ie.SrcInterfaceAccess),
			ie.NewFTEID(0x01, b.teid, net.ParseIP(b.n3Address), nil, 0),
			ie.NewSDFFilter(b.sdfFilter, "", "", "", 1),
		),
		ie.NewOuterHeaderRemoval(0, 0),
		ie.NewFARID(b.farID),
		ie.NewQERID(b.sessQerID),
		ie.NewQERID(b.appQerID),
	)

}
