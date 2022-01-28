// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package session

type IEMethod uint8

// Definitions for session rules

const (
	Create IEMethod = iota
	Update
	Delete

	ActionForward uint8 = 0x2
	ActionDrop    uint8 = 0x1
	ActionBuffer  uint8 = 0x4
	ActionNotify  uint8 = 0x8
)
