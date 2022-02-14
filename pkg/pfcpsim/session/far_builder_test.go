package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				method:      Create,
				applyAction: []uint8{ActionDrop},
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
				farID:       2,
				applyAction: []uint8{ActionDrop},
				teid:        100,
			},
			description: "Invalid FAR: Providing TEID without DownlinkIP",
		},
		{
			input: NewFARBuilder().WithMethod(Create).
				WithID(1).
				WithAction(ActionForward).
				WithDownlinkIP("10.0.0.1"),
			expected: &farBuilder{
				farID:       1,
				applyAction: []uint8{ActionForward},
				downlinkIP:  "10.0.0.1",
			},
			description: "Invalid FAR: Providing DownlinkIP without TEID",
		},
		{
			input: NewFARBuilder().WithMethod(Create).
				WithID(1).
				WithAction(ActionForward).
				WithAction(ActionDrop).
				WithAction(ActionBuffer).
				WithAction(ActionNotify).
				WithDownlinkIP("10.0.0.1").
				WithTEID(100),
			expected: &farBuilder{
				farID:       1,
				method:      Create,
				applyAction: []uint8{ActionForward, ActionDrop, ActionBuffer, ActionNotify},
				downlinkIP:  "10.0.0.1",
				teid:        100,
			},
			description: "Invalid FAR: Providing both Forward and Drop actions",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.Panics(t, func() { scenario.input.BuildFAR() })
			assert.Equal(t, scenario.input, scenario.expected)
		})
	}
}
