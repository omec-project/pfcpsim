package session

import (
	"github.com/wmnsk/go-pfcp/ie"
)

type farBuilder struct {
	id          uint32
	applyAction uint8
	method      IEMethod
}

type downlinkFARBuilder struct {
	farBuilder
	teid       uint32
	downlinkIP string
}

func NewFARBuilder() *farBuilder {
	return &farBuilder{}
}

func (b *farBuilder) WithID(id uint32) *farBuilder {
	b.id = id
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

func (b *farBuilder) BuildUplinkFAR() *ie.IE {
	createFunc := ie.NewCreateFAR
	if b.method == Update {
		createFunc = ie.NewUpdateFAR
	}

	return createFunc(
		ie.NewFARID(b.id),
		ie.NewApplyAction(b.applyAction),
		ie.NewForwardingParameters(
			ie.NewDestinationInterface(ie.DstInterfaceCore),
		),
	)
}

func (b *farBuilder) WithTEID(teid uint32) *downlinkFARBuilder {
	return &downlinkFARBuilder{
		farBuilder: *b,
		teid:       teid,
	}
}

func (b *farBuilder) WithDownlinkIP(downlinkIP string) *downlinkFARBuilder {
	return &downlinkFARBuilder{
		farBuilder: *b,
		downlinkIP: downlinkIP,
	}
}

func (b *downlinkFARBuilder) BuildDownlinkFAR() *ie.IE {
	createFunc := ie.NewCreateFAR
	if b.method == Update {
		createFunc = ie.NewUpdateFAR
	}

	return createFunc(
		ie.NewFARID(b.id),
		ie.NewApplyAction(b.applyAction),
		ie.NewUpdateForwardingParameters(
			ie.NewDestinationInterface(ie.DstInterfaceAccess),
			// FIXME desc 0x100?
			ie.NewOuterHeaderCreation(0x100, b.teid, b.downlinkIP, "", 0, 0, 0),
		),
	)
}
