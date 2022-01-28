package session

import (
	"github.com/wmnsk/go-pfcp/ie"
)

type farBuilder struct {
	farID         uint32
	applyAction   uint8
	method        IEMethod
	teid          uint32
	downlinkIP    string
	isDownlinkFAR bool
}

// NewFARBuilder returns a farBuilder.
func NewFARBuilder() *farBuilder {
	return &farBuilder{
		farID:       uint32(1),
		method:      Create,
		applyAction: ActionDrop,
	}
}

func (b *farBuilder) WithID(id uint32) *farBuilder {
	b.farID = id
	return b
}

func (b *farBuilder) WithMethod(method IEMethod) *farBuilder {
	b.method = method
	return b
}

func (b *farBuilder) WithAction(action uint8) *farBuilder {
	b.applyAction = action
	return b
}

func (b *farBuilder) WithTEID(teid uint32) *farBuilder {
	b.teid = teid
	return b
}

func (b *farBuilder) WithDownlinkIP(downlinkIP string) *farBuilder {
	b.downlinkIP = downlinkIP
	return b
}

func (b *farBuilder) MarkDownlink() *farBuilder {
	b.isDownlinkFAR = true
	return b
}

// BuildFAR returns by default an UplinkFAR.
// Returns a DownlinkFAR if MarkDownlink was invoked.
func (b *farBuilder) BuildFAR() *ie.IE {
	createFunc := ie.NewCreateFAR
	if b.method == Update {
		createFunc = ie.NewUpdateFAR
	}

	if b.isDownlinkFAR {
		return createFunc(
			ie.NewFARID(b.farID),
			ie.NewApplyAction(b.applyAction),
			ie.NewUpdateForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				// FIXME desc 0x100?
				ie.NewOuterHeaderCreation(0x100, b.teid, b.downlinkIP, "", 0, 0, 0),
			),
		)
	}

	// Uplink
	return createFunc(
		ie.NewFARID(b.farID),
		ie.NewApplyAction(b.applyAction),
		ie.NewForwardingParameters(
			ie.NewDestinationInterface(ie.DstInterfaceCore),
		),
	)
}
