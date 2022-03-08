// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Open Networking Foundation

package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wmnsk/go-pfcp/ie"
)

func TestQERBuilderShouldPanic(t *testing.T) {
	type testCase struct {
		input       *qerBuilder
		expected    *qerBuilder
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewQERBuilder().
				WithMethod(Create).
				WithQFI(1),
			expected: &qerBuilder{
				method: Create,
				qfi:    1,
			},
			description: "Invalid QER: No ID provided",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.Panics(t, func() { scenario.input.Build() })
			assert.Equal(t, scenario.input, scenario.expected)
		})
	}
}

func TestQERBuilder(t *testing.T) {
	type testCase struct {
		input       *qerBuilder
		expected    *ie.IE
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewQERBuilder().
				WithID(1).
				WithMethod(Create).
				WithQFI(2),
			expected: ie.NewCreateQER(
				ie.NewQERID(1),
				ie.NewQFI(2),
				ie.NewGateStatus(0, 0),
			),
			description: "Valid Create QER",
		},
		{
			input: NewQERBuilder().
				WithID(1).
				WithMethod(Update).
				WithQFI(2).
				WithDownlinkMBR(0).
				WithUplinkMBR(0),
			expected: ie.NewUpdateQER(
				ie.NewQERID(1),
				ie.NewQFI(2),
				ie.NewGateStatus(0, 0),
				ie.NewMBR(0, 0),
			),
			description: "Valid Update QER",
		},
		{
			input: NewQERBuilder().
				WithID(1).
				WithMethod(Delete).
				WithQFI(2).
				WithDownlinkGBR(0).
				WithUplinkGBR(0),
			expected: ie.NewRemoveQER(
				ie.NewCreateQER(
					ie.NewQERID(1),
					ie.NewQFI(2),
					ie.NewGateStatus(0, 0),
					ie.NewGBR(0, 0),
				),
			),
			description: "Valid Delete QER",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.NotPanics(t, func() { _ = scenario.input.Build() })
			assert.Equal(t, scenario.input.Build(), scenario.expected)
		})
	}
}
