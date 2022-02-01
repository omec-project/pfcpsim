package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wmnsk/go-pfcp/ie"
)

func TestFARBuilderShouldPanic(t *testing.T) {
	type testCase struct {
		input *farBuilder
		expected *ie.IE
		description string
	}

	for _, scenario := range []testCase{
		{
			input: NewFARBuilder().
				WithMethod(Create).
				WithAction(ActionDrop),
			expected: nil,
			description: "Invalid FAR: No ID provided",
		},
		{
			input: NewFARBuilder().
				WithMethod(Create).
			WithID(2).
			WithAction(ActionDrop).
			WithTEID(100),

			expected: nil,
			description: "Invalid FAR: Providing TEID without DownlinkIP",
		},
		{
			input: NewFARBuilder().WithMethod(Create).
				WithID(1).
				WithAction(ActionDrop).
				WithDownlinkIP("10.0.0.1"),
			expected: nil,
			description: "Invalid FAR: Providing DownlinkIP without TEID",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.Panics(t, func() {scenario.input.BuildFAR()} )
		})
	}
}
