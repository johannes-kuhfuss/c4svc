package providers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFile = "D:\\Data\\code\\github.com\\johannes-kuhfuss\\c4\\media\\TestBild.tif"
)

func TestProcessFileOpenError(t *testing.T) {
	c4Id, err := C4Provider.ProcessFile("", "", false)
	assert.NotNil(t, err)
	assert.EqualValues(t, "", c4Id)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "Source file not found or could not be read", err.Message())
}
