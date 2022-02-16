package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wmnsk/go-pfcp/ie"
)

func TestFARBuilderShouldPanic(t *testing.T) {
	type testCase struct {
		input       *farBuilder
		expected    *farBuilder
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewFARBuilder().
				WithMethod(Create).
				WithAction(ActionDrop),
			expected: &farBuilder{
				method:  Create,
				actions: []uint8{ActionDrop},
			},
			description: "Invalid FAR: No ID provided",
		},
		{
			input: NewFARBuilder().
				WithMethod(Create).
				WithID(2).
				WithAction(ActionDrop).
				WithTEID(100),

			expected: &farBuilder{
				farID:   2,
				method:  Create,
				actions: []uint8{ActionDrop},
				teid:    100,
			},
			description: "Invalid FAR: Providing TEID without DownlinkIP",
		},
		{
			input: NewFARBuilder().WithMethod(Create).
				WithID(1).
				WithAction(ActionForward).
				WithDownlinkIP("10.0.0.1"),
			expected: &farBuilder{
				farID:      1,
				method:     Create,
				actions:    []uint8{ActionForward},
				downlinkIP: "10.0.0.1",
			},
			description: "Invalid FAR: Providing DownlinkIP without TEID",
		},
		{
			input: NewFARBuilder().WithMethod(Create).
				WithID(1).
				WithAction(ActionForward).
				WithAction(ActionDrop).
				WithDownlinkIP("10.0.0.1").
				WithTEID(100),
			expected: &farBuilder{
				farID:      1,
				method:     Create,
				actions:    []uint8{ActionForward, ActionDrop},
				downlinkIP: "10.0.0.1",
				teid:       100,
			},
			description: "Invalid FAR: Providing both forward and drop actions",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.Panics(t, func() { scenario.input.BuildFAR() })
			assert.Equal(t, scenario.input, scenario.expected)
		})
	}
}

func TestFARBuilder(t *testing.T) {
	type testCase struct {
		input       *farBuilder
		expected    *ie.IE
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewFARBuilder().
				WithID(1).
				WithMethod(Create).
				WithAction(ActionForward).
				WithDstInterface(ie.DstInterfaceAccess),
			expected: ie.NewCreateFAR(
				ie.NewFARID(1),
				ie.NewForwardingParameters(
					ie.NewDestinationInterface(ie.DstInterfaceAccess),
				),
				ie.NewApplyAction(ActionForward),
			),
			description: "Valid FAR",
		},
		{
			input: NewFARBuilder().
				WithID(1).
				WithMethod(Create).
				WithAction(ActionDrop).
				WithAction(ActionBuffer).
				WithDstInterface(ie.DstInterfaceAccess).
				WithTEID(12).
				WithDownlinkIP("10.0.0.1"),
			expected: ie.NewCreateFAR(
				ie.NewFARID(1),
				ie.NewForwardingParameters(
					ie.NewDestinationInterface(ie.DstInterfaceAccess),
					ie.NewOuterHeaderCreation(S_TAG, 12, "10.0.0.1", "", 0, 0, 0),
				),
				ie.NewApplyAction(ActionDrop),
				ie.NewApplyAction(ActionBuffer),
			),
			description: "Valid FAR with 2 actions",
		},
		{
			input: NewFARBuilder().
				WithID(1).
				WithMethod(Create).
				WithAction(ActionForward).
				WithAction(ActionBuffer).
				WithAction(ActionNotify).
				WithDstInterface(ie.DstInterfaceAccess),
			expected: ie.NewCreateFAR(
				ie.NewFARID(1),
				ie.NewForwardingParameters(
					ie.NewDestinationInterface(ie.DstInterfaceAccess),
				),
				ie.NewApplyAction(ActionForward),
				ie.NewApplyAction(ActionBuffer),
				ie.NewApplyAction(ActionNotify),
			),
			description: "Valid FAR 3 actions",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			require.Equal(t, scenario.expected, scenario.input.BuildFAR())
		})
	}
}
