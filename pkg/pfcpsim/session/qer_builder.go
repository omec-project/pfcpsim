// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package session

import (
	"github.com/omec-project/pfcpsim/logger"
	"github.com/wmnsk/go-pfcp/ie"
)

type qerBuilder struct {
	method     IEMethod
	qerID      uint32
	qfi        uint8
	isMbrSet   bool
	ulMbr      uint64
	dlMbr      uint64
	isGbrSet   bool
	ulGbr      uint64
	dlGbr      uint64
	gateStatus uint8

	isIDSet bool
}

const (
	QerNoFuzz          = 0
	QerWithQFI         = 1
	QerWithUplinkMBR   = 2
	QerWithUplinkGBR   = 3
	QerWithDownlinkMBR = 4
	QerWithDownlinkGBR = 5
	QerWithGateStatus  = 6
	QerMax             = 7
)

func NewQERBuilder() *qerBuilder {
	return &qerBuilder{}
}

func (b *qerBuilder) FuzzIE(ieType int, arg uint) *qerBuilder {
	switch ieType {
	case QerWithQFI:
		logger.PfcpsimLog.Infoln("QerWithQFI")
		return b.WithQFI(uint8(arg))
	case QerWithUplinkMBR:
		logger.PfcpsimLog.Infoln("QerWithUplinkMBR")
		return b.WithUplinkMBR(uint64(arg))
	case QerWithUplinkGBR:
		logger.PfcpsimLog.Infoln("QerWithUplinkGBR")
		return b.WithUplinkGBR(uint64(arg))
	case QerWithDownlinkMBR:
		logger.PfcpsimLog.Infoln("QerWithDownlinkMBR")
		return b.WithDownlinkMBR(uint64(arg))
	case QerWithDownlinkGBR:
		logger.PfcpsimLog.Infoln("QerWithDownlinkGBR")
		return b.WithDownlinkGBR(uint64(arg))
	case QerWithGateStatus:
		logger.PfcpsimLog.Infoln("QerWithGateStatus")
		return b.WithGateStatus(uint8(arg))
	default:
	}

	return b
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

func (b *qerBuilder) WithGateStatus(status uint8) *qerBuilder {
	b.gateStatus = status

	return b
}

func (b *qerBuilder) validate() {
	if !b.isIDSet {
		logger.PfcpsimLog.Panicln("tried to build a QER without setting the QER ID")
	}
}

func (b *qerBuilder) WithMethod(method IEMethod) *qerBuilder {
	b.method = method
	return b
}

func (b *qerBuilder) Build() *ie.IE {
	if doCheck {
		b.validate()
	}

	createFunc := ie.NewCreateQER
	if b.method == Update {
		createFunc = ie.NewUpdateQER
	}

	gate := ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusOpen)
	if b.gateStatus == ie.GateStatusClosed {
		gate = ie.NewGateStatus(ie.GateStatusClosed, ie.GateStatusClosed)
	}

	qer := createFunc(
		ie.NewQERID(b.qerID),
		ie.NewQFI(b.qfi),
		gate,
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
