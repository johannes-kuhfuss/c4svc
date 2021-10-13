package date_utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	assert.EqualValues(t, time.RFC3339, apiDateLayout)
}

func TestGetNowUtcString(t *testing.T) {
	now := GetNowUtcString()
	date, err := time.Parse(time.RFC3339, now)
	assert.NotNil(t, date)
	assert.Nil(t, err)
	assert.EqualValues(t, date.Format(time.RFC3339), now)
}
