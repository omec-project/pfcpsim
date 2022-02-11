/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package session

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wmnsk/go-pfcp/ie"
)

func TestPDRBuilderShouldPanic(t *testing.T) {
	type testCase struct {
		input       *pdrBuilder
		expected    *pdrBuilder
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewPDRBuilder().
				WithMethod(Create).
				MarkAsDownlink(),
			expected: &pdrBuilder{
				method:    Create,
				direction: downlink,
				qerIDs:    make([]*ie.IE, 0),
			},
			description: "Invalid Downlink PDR: No ID provided",
		},
		{
			input: NewPDRBuilder().
				WithMethod(Create).
				WithTEID(100).
				MarkAsDownlink(),
			expected: &pdrBuilder{
				method:    Create,
				direction: downlink,
				teid:      100,
				qerIDs:    make([]*ie.IE, 0),
			},
			description: "Invalid Downlink PDR: Partial parameters provided",
		},
		{
			input: NewPDRBuilder().
				WithMethod(Create).
				WithUEAddress("10.0.0.1").
				MarkAsDownlink(),
			expected: &pdrBuilder{
				method:    Create,
				direction: downlink,
				ueAddress: "10.0.0.1",
				qerIDs:    make([]*ie.IE, 0),
			},
			description: "Invalid Downlink PDR: Partial parameters provided",
		},
		{
			input: NewPDRBuilder().
				WithMethod(Create).
				WithUEAddress("10.0.0.1").
				WithTEID(100),
			expected: &pdrBuilder{
				method:    Create,
				ueAddress: "10.0.0.1",
				teid:      100,
				qerIDs:    make([]*ie.IE, 0),
			},
			description: "Invalid Downlink PDR: build without MarkAsDownlink",
		},
		{
			input: NewPDRBuilder().
				WithMethod(Create).
				WithUEAddress("10.0.0.1").
				WithTEID(100).
				MarkAsUplink(),
			expected: &pdrBuilder{
				method:    Create,
				ueAddress: "10.0.0.1",
				direction: uplink,
				teid:      100,
				qerIDs:    make([]*ie.IE, 0),
			},
			description: "Invalid Downlink PDR: marked as uplink passing downlink parameters",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.Panics(t, func() { scenario.input.BuildPDR() })
			assert.Equal(t, scenario.input, scenario.expected)
		})
	}
}

func TestPDRBuilder(t *testing.T) {
	type testCase struct {
		input       *pdrBuilder
		expected    *ie.IE
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewPDRBuilder().
				WithID(1).
				WithPrecedence(2).
				WithTEID(100).
				WithMethod(Create).
				WithN3Address("192.168.0.1").
				WithFARID(3).
				AddQERID(4).
				WithSDFFilter("permit ip any to assigned").
				MarkAsUplink(),
			expected: ie.NewCreatePDR(
				ie.NewPDRID(1),
				ie.NewPrecedence(2),
				ie.NewPDI(
					ie.NewSourceInterface(ie.SrcInterfaceAccess),
					ie.NewFTEID(0x01, 100, net.ParseIP("192.168.0.1"), nil, 0),
					ie.NewSDFFilter("permit ip any to assigned", "", "", "", 1),
				),
				ie.NewOuterHeaderRemoval(0, 0),
				ie.NewFARID(3),
				ie.NewQERID(4),
			),
			description: "Valid Create Uplink PDR",
		},
		{
			input: NewPDRBuilder().
				WithID(1).
				WithPrecedence(2).
				WithTEID(100).
				WithUEAddress("172.16.0.1").
				WithMethod(Update).
				WithFARID(3).
				AddQERID(4).
				WithSDFFilter("permit ip any to assigned").
				MarkAsDownlink(),
			expected: ie.NewUpdatePDR(
				ie.NewPDRID(1),
				ie.NewPrecedence(2),
				ie.NewPDI(
					ie.NewSourceInterface(ie.SrcInterfaceCore),
					ie.NewUEIPAddress(0x2, "172.16.0.1", "", 0, 0),
					ie.NewSDFFilter("permit ip any to assigned", "", "", "", 1),
				),
				ie.NewFARID(3),
				ie.NewQERID(4),
			),
			description: "Valid Update Downlink PDR",
		},
		{
			input: NewPDRBuilder().
				WithID(1).
				WithPrecedence(2).
				WithTEID(100).
				WithUEAddress("172.16.0.1").
				WithMethod(Delete).
				WithFARID(3).
				AddQERID(4).
				WithSDFFilter("permit ip any to assigned").
				MarkAsDownlink(),
			expected: ie.NewRemovePDR(
				ie.NewCreatePDR(
					ie.NewPDRID(1),
					ie.NewPrecedence(2),
					ie.NewPDI(
						ie.NewSourceInterface(ie.SrcInterfaceCore),
						ie.NewUEIPAddress(0x2, "172.16.0.1", "", 0, 0),
						ie.NewSDFFilter("permit ip any to assigned", "", "", "", 1),
					),
					ie.NewFARID(3),
					ie.NewQERID(4),
				),
			),
			description: "Valid Delete Downlink PDR",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.NotPanics(t, func() { scenario.input.BuildPDR() })
			assert.Equal(t, scenario.input.BuildPDR(), scenario.expected)
		})
	}
}
