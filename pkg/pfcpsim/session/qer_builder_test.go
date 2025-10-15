// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package session

import (
	"reflect"
	"testing"

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
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected Build() to panic, but it didn't")
				}
			}()
			scenario.input.Build()

			if !reflect.DeepEqual(scenario.input, scenario.expected) {
				t.Errorf("QER builder mismatch. got = %+v, want = %+v", scenario.input, scenario.expected)
			}
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
			description: "Valid Create QER with gate open",
		},
		{
			input: NewQERBuilder().
				WithID(1).
				WithMethod(Create).
				WithQFI(2).
				WithGateStatus(ie.GateStatusClosed),
			expected: ie.NewCreateQER(
				ie.NewQERID(1),
				ie.NewQFI(2),
				ie.NewGateStatus(1, 1),
			),
			description: "Valid Create QER with Gate closed",
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
			var result *ie.IE
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Build() panicked unexpectedly: %v", r)
					}
				}()
				result = scenario.input.Build()
			}()

			if !reflect.DeepEqual(result, scenario.expected) {
				t.Errorf("QER build result mismatch. got = %+v, want = %+v", result, scenario.expected)
			}
		})
	}
}
