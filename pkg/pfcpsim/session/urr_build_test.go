// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-present Ian Chen <ychen.cs10@nycu.edu.tw>

package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			assert.Panics(t, func() { scenario.input.Build() })
			assert.Equal(t, scenario.input, scenario.expected)
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
				ie.NewReportingTriggers(RPT_TRIG_PERIO),
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
					ie.NewReportingTriggers(RPT_TRIG_PERIO),
				),
			),
			description: "Valid Delete URR",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.NotPanics(t, func() { _ = scenario.input.Build() })
			assert.Equal(t, scenario.expected, scenario.input.Build())
		})
	}
}
