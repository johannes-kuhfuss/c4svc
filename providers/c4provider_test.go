package providers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFile = "D:\\Data\\code\\github.com\\johannes-kuhfuss\\c4srv\\media\\TestBild.tif"
)

func TestProcessFileOpenError(t *testing.T) {
	c4Id, err := C4Provider.ProcessFile("")
	assert.NotNil(t, err)
	assert.Nil(t, c4Id)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "Source file not found or could not be read", err.Message())
}

func TestProcessFileNoError(t *testing.T) {
	c4Id, err := C4Provider.ProcessFile(testFile)
	assert.NotNil(t, c4Id)
	assert.Nil(t, err)
	assert.EqualValues(t, "c42FTpMRKrEEL6sgwRVRfxbzYDsYZe4VgsNVC7D6Jkqz8ABjsSAybKLYwPLGSJexGkJ9qt3aR8sMAjZ8fhKd7GfQsB", *c4Id)
}
