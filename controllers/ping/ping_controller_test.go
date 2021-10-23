package ping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "pong", pong)
}

func TestPing(t *testing.T) {
	// To do
}
