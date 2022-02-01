package session

import (
	"github.com/wmnsk/go-pfcp/ie"
)

type farBuilder struct {
	farID       uint32
	applyAction uint8
	method      IEMethod
	teid        uint32
	downlinkIP  string
	direction   direction
}

// NewFARBuilder returns a farBuilder.
func NewFARBuilder() *farBuilder {
	return &farBuilder{}
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

func (b *farBuilder) MarkAsDownlink() *farBuilder {
	b.direction = downlink
	return b
}

func (b *farBuilder) MarkAsUplink() *farBuilder {
	b.direction = uplink
	return b
}

func (b *farBuilder) validate() {
	if b.direction == notSet {
		panic("Tried building a FAR without marking it as uplink or downlink")
	}
}

func newRemoveFAR(far *ie.IE) *ie.IE {
	return ie.NewRemoveFAR(far)
}

// BuildFAR returns by default a downlinkFAR if MarkAsUplink was invoked.
// Returns a DownlinkFAR if MarkAsDownlink was invoked.
func (b *farBuilder) BuildFAR() *ie.IE {
	b.validate()

	createFunc := ie.NewCreateFAR
	if b.method == Update {
		createFunc = ie.NewUpdateFAR
	}

	if b.direction == downlink {
		far := createFunc(
			ie.NewFARID(b.farID),
			ie.NewApplyAction(b.applyAction),
			ie.NewUpdateForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				// FIXME desc 0x100?
				ie.NewOuterHeaderCreation(0x100, b.teid, b.downlinkIP, "", 0, 0, 0),
			),
		)
		if b.method == Delete {
			return newRemoveFAR(far)
		}

		return far
	}

	// Uplink
	far := createFunc(
		ie.NewFARID(b.farID),
		ie.NewApplyAction(b.applyAction),
		ie.NewForwardingParameters(
			ie.NewDestinationInterface(ie.DstInterfaceCore),
		),
	)

	if b.method == Delete {
		return newRemoveFAR(far)
	}

	return far
}
