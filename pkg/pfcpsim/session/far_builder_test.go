package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFARBuilderShouldPanic(t *testing.T) {
	assert.Panics(t, func() {
		NewFARBuilder().
			WithMethod(Create).
			WithAction(ActionDrop).
			BuildFAR()
	})

	// Providing partial parameters (DownlinkIP without providing TEID)
	assert.Panics(t, func() {
		NewFARBuilder().WithMethod(Create).
			WithID(1).
			WithAction(ActionDrop).
			WithDownlinkIP("10.0.0.1").
			BuildFAR()
	})

	// Providing partial parameters (TEID without providing DownlinkIP)
	assert.Panics(t, func() {
		NewFARBuilder().WithMethod(Create).
			WithID(2).
			WithAction(ActionDrop).
			WithTEID(100).
			BuildFAR()
	})
}
