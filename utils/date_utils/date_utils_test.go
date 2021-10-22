package date_utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	assert.EqualValues(t, time.RFC3339, ApiDateLayout)
}

func TestGetNowUtcString(t *testing.T) {
	now := GetNowUtcString()
	date, err := time.Parse(ApiDateLayout, now)
	assert.NotNil(t, date)
	assert.Nil(t, err)
	assert.EqualValues(t, date.Format(ApiDateLayout), now)
}
