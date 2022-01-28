// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package session

import (
	"net"

	"github.com/wmnsk/go-pfcp/ie"
)

type IEMethod uint8

const (
	Create IEMethod = iota
	Update
	Delete
)

const (
	ActionForward uint8 = 0x2
	ActionDrop    uint8 = 0x1
	ActionBuffer  uint8 = 0x4
	ActionNotify  uint8 = 0x8
)

const (
	dummyPrecedence = 100
)

// TODO: use builder pattern to create PDR IE
func NewUplinkPDR(method IEMethod, id uint16, teid uint32, n3address string,
	farID uint32, sessQerID uint32, appQerID uint32) *ie.IE {
	createFunc := ie.NewCreatePDR
	if method == Update {
		createFunc = ie.NewUpdatePDR
	}

	return createFunc(
		ie.NewPDRID(id),
		ie.NewPrecedence(dummyPrecedence),
		ie.NewPDI(
			ie.NewSourceInterface(ie.SrcInterfaceAccess),
			ie.NewFTEID(0x01, teid, net.ParseIP(n3address), nil, 0),
			ie.NewSDFFilter("permit out ip from any to assigned", "", "", "", 1),
		),
		ie.NewOuterHeaderRemoval(0, 0),
		ie.NewFARID(farID),
		ie.NewQERID(appQerID),
		ie.NewQERID(sessQerID),
	)
}

func NewDownlinkPDR(method IEMethod, id uint16, ueAddress string,
	farID uint32, sessQerID uint32, appQerID uint32) *ie.IE {
	createFunc := ie.NewCreatePDR
	if method == Update {
		createFunc = ie.NewUpdatePDR
	}

	return createFunc(
		ie.NewPDRID(id),
		ie.NewPrecedence(dummyPrecedence),
		ie.NewPDI(
			ie.NewSourceInterface(ie.SrcInterfaceCore),
			ie.NewUEIPAddress(0x2, ueAddress, "", 0, 0),
			ie.NewSDFFilter("permit out ip from any to assigned", "", "", "", 1),
		),
		ie.NewFARID(farID),
		ie.NewQERID(appQerID),
		ie.NewQERID(sessQerID),
	)
}
