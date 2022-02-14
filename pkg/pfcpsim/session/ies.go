// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package session

type IEMethod uint8

// Definitions for session rules

type direction int

const (
	notSet direction = iota
	uplink
	downlink
)

const (
	Create IEMethod = iota
	Update
	Delete

	ActionForward uint8 = 0x2
	ActionDrop    uint8 = 0x1
	ActionBuffer  uint8 = 0x4
	ActionNotify  uint8 = 0x8

	S_TAG = 0x100 // 8th bit set. Used in FAR's Outer Header Creation as description. Refer to Table 8.2.56-1
)
