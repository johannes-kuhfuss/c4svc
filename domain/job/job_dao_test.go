package domain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNotFound(t *testing.T) {
	user, err := JobDao.Get("X")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "job with Id X does not exist", err.Message())
}

func TestGetNoError(t *testing.T) {
	user, err := JobDao.Get("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, "1zXgBZNnBG1msmF1ARQK9ZphbbO", user.Id)
}
