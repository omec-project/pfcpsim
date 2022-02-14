package session

import "github.com/wmnsk/go-pfcp/ie"

type qerBuilder struct {
	method   IEMethod
	qerID    uint32
	qfi      uint8
	isMbrSet bool
	ulMbr    uint64
	dlMbr    uint64
	isGbrSet bool
	ulGbr    uint64
	dlGbr    uint64

	isIDSet bool
}

func NewQERBuilder() *qerBuilder {
	return &qerBuilder{}
}

func (b *qerBuilder) WithID(id uint32) *qerBuilder {
	// Used to avoid using 0 as default value. It makes sure that WithID was invoked.
	b.isIDSet = true
	b.qerID = id

	return b
}

func (b *qerBuilder) WithQFI(qfi uint8) *qerBuilder {
	b.qfi = qfi
	return b
}

func (b *qerBuilder) WithUplinkMBR(ulMbr uint64) *qerBuilder {
	b.isMbrSet = true
	b.ulMbr = ulMbr

	return b
}

func (b *qerBuilder) WithUplinkGBR(ulGbr uint64) *qerBuilder {
	b.isGbrSet = true
	b.ulGbr = ulGbr

	return b
}

func (b *qerBuilder) WithDownlinkMBR(dlMbr uint64) *qerBuilder {
	b.isMbrSet = true
	b.dlMbr = dlMbr

	return b
}

func (b *qerBuilder) WithDownlinkGBR(dlGbr uint64) *qerBuilder {
	b.isGbrSet = true
	b.dlGbr = dlGbr

	return b
}

func (b *qerBuilder) validate() {
	if !b.isIDSet {
		panic("Tried to build a QER without setting the QER ID")
	}
}

func (b *qerBuilder) WithMethod(method IEMethod) *qerBuilder {
	b.method = method
	return b
}

func (b *qerBuilder) Build() *ie.IE {
	b.validate()

	createFunc := ie.NewCreateQER
	if b.method == Update {
		createFunc = ie.NewUpdateQER
	}

	qer := createFunc(
		ie.NewQERID(b.qerID),
		ie.NewQFI(b.qfi),
		// FIXME: we don't support gating, always OPEN
		ie.NewGateStatus(0, 0),
	)

	if b.isMbrSet {
		qer.Add(ie.NewMBR(b.ulMbr, b.dlMbr))
	}

	if b.isGbrSet {
		qer.Add(ie.NewGBR(b.ulGbr, b.dlGbr))
	}

	if b.method == Delete {
		return ie.NewRemoveQER(qer)
	}

	return qer
}
