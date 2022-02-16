package session

import (
	"github.com/wmnsk/go-pfcp/ie"
)

type farBuilder struct {
	farID        uint32
	actions      []uint8
	method       IEMethod
	teid         uint32
	downlinkIP   string
	dstInterface uint8

	zeroBasedOuterHeader bool
	isInterfaceSet       bool
}

// NewFARBuilder returns a farBuilder.
func NewFARBuilder() *farBuilder {
	return &farBuilder{}
}

func (b *farBuilder) WithID(id uint32) *farBuilder {
	b.farID = id
	return b
}

func (b *farBuilder) WithZeroBasedOuterHeaderCreation() *farBuilder {
	b.zeroBasedOuterHeader = true
	return b
}

func (b *farBuilder) WithMethod(method IEMethod) *farBuilder {
	b.method = method
	return b
}

func (b *farBuilder) WithAction(action ...uint8) *farBuilder {
	b.actions = append(b.actions, action...)
	return b
}

func (b *farBuilder) WithTEID(teid uint32) *farBuilder {
	b.teid = teid
	return b
}

func (b *farBuilder) WithDstInterface(iFace uint8) *farBuilder {
	b.isInterfaceSet = true
	b.dstInterface = iFace

	return b
}

func (b *farBuilder) WithDownlinkIP(downlinkIP string) *farBuilder {
	b.downlinkIP = downlinkIP
	return b
}

func (b *farBuilder) validate() {
	if b.farID == 0 {
		panic("Tried building FAR without setting FAR ID")
	}

	if len(b.actions) == 0 {
		panic("Tried building FAR without setting at least one action")
	}

	if len(b.actions) != 0 {
		if contains(b.actions, ActionDrop) && contains(b.actions, ActionForward) {
			panic("Tried building FAR with both forward and drop actions")
		}
	}

	if !b.isInterfaceSet {
		panic("Tried building FAR without setting a destination interface")
	}

	if (b.downlinkIP != "" && b.teid == 0 || b.downlinkIP == "" && b.teid != 0) && !b.zeroBasedOuterHeader {
		panic("Tried building FAR providing only partial parameters. Check downlink IP or TEID")
	}
}

// BuildFAR returns by default a downlinkFAR if MarkAsUplink was invoked.
// Returns a DownlinkFAR if MarkAsDownlink was invoked.
func (b *farBuilder) BuildFAR() *ie.IE {
	b.validate()

	fwdParams := ie.NewForwardingParameters(
		ie.NewDestinationInterface(b.dstInterface),
	)

	createFunc := ie.NewCreateFAR
	if b.method == Update {
		createFunc = ie.NewUpdateFAR
		fwdParams = ie.NewUpdateForwardingParameters(
			ie.NewDestinationInterface(b.dstInterface),
		)
	}

	if b.zeroBasedOuterHeader {
		fwdParams.Add(ie.NewOuterHeaderCreation(S_TAG, 0, "0.0.0.0", "", 0, 0, 0))
	} else if b.teid != 0 && b.downlinkIP != "" {
		// TEID and DownlinkIP are provided
		fwdParams.Add(ie.NewOuterHeaderCreation(S_TAG, b.teid, b.downlinkIP, "", 0, 0, 0))
	}

	far := createFunc(
		ie.NewFARID(b.farID),
		fwdParams,
	)

	for _, action := range b.actions {
		far.Add(ie.NewApplyAction(action))
	}

	if b.method == Delete {
		return ie.NewRemoveFAR(far)
	}

	return far
}
