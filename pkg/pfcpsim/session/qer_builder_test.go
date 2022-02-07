/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		{
			input: NewQERBuilder().
				WithMethod(Create).
				WithID(2),
			expected: &qerBuilder{
				method:  Create,
				qerID:   2,
				isIDSet: true,
			},
			description: "Invalid QER: No QFI provided",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.Panics(t, func() { scenario.input.Build() })
			assert.Equal(t, scenario.input, scenario.expected)
		})
	}
}
