// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-present Ian Chen <ychen.cs10@nycu.edu.tw>

package session

import (
	"reflect"
	"testing"
	"time"

	"github.com/wmnsk/go-pfcp/ie"
)

func TestURRBuilderShouldPanic(t *testing.T) {
	type testCase struct {
		input       *urrBuilder
		expected    *urrBuilder
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewURRBuilder().
				WithID(1).
				WithMethod(Create).
				WithMeasurementMethod(0, 1, 0).
				WithMeasurementPeriod(0).
				WithMeasurementInfo(17).
				WithReportingTrigger(ReportingTrigger{
					Flags: RPT_TRIG_PERIO,
				}),
			expected: &urrBuilder{
				method:          Create,
				urrID:           1,
				measurementInfo: 17,
				measurementMethod: &measurementMethodParams{
					event: 0,
					volum: 1,
					durat: 0,
				},
				measurementPeriod: 0,
				rptTrig: ReportingTrigger{
					Flags: RPT_TRIG_PERIO,
				},
			},
			description: "Invalid URR: measurementInfo is invalid",
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

func TestURRBuilder(t *testing.T) {
	type testCase struct {
		input       *urrBuilder
		expected    *ie.IE
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewURRBuilder().
				WithID(1).
				WithMethod(Create).
				WithMeasurementMethod(0, 1, 0).
				WithMeasurementPeriod(time.Second).
				WithReportingTrigger(ReportingTrigger{
					Flags: RPT_TRIG_PERIO,
				}),
			expected: ie.NewCreateURR(
				ie.NewURRID(1),
				ie.NewMeasurementMethod(0, 1, 0),
				ie.NewMeasurementPeriod(time.Second),
				NewRptTrig(ReportingTrigger{
					Flags: RPT_TRIG_PERIO,
				}),
			),
			description: "Valid Create URR with reporting trigger",
		},
		{
			input: NewURRBuilder().
				WithID(1).
				WithMethod(Update).
				WithMeasurementPeriod(2 * time.Second),
			expected: ie.NewUpdateURR(
				ie.NewURRID(1),
				ie.NewMeasurementPeriod(2*time.Second),
			),
			description: "Valid Update URR",
		},
		{
			input: NewURRBuilder().
				WithID(1).
				WithMethod(Delete).
				WithMeasurementMethod(0, 1, 0).
				WithMeasurementPeriod(1).
				WithReportingTrigger(ReportingTrigger{
					Flags: RPT_TRIG_PERIO,
				}),
			expected: ie.NewRemoveURR(
				ie.NewCreateURR(
					ie.NewURRID(1),
					ie.NewMeasurementMethod(0, 1, 0),
					ie.NewMeasurementPeriod(2),
					NewRptTrig(ReportingTrigger{
						Flags: RPT_TRIG_PERIO,
					}),
				),
			),
			description: "Valid Delete URR",
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
