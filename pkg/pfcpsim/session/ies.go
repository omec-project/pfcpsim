// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package session

type IEMethod uint8

// Definitions for session rules

type direction int

const (
	notSet direction = iota
	uplink
	downlink

	Create IEMethod = iota
	Update
	Delete

	ActionForward uint8 = 0x2
	ActionDrop    uint8 = 0x1
	ActionBuffer  uint8 = 0x4
	ActionNotify  uint8 = 0x8

	S_TAG = 0x100 // Refer to table 8.2.56-1 in PFCP specs Release 16
)
