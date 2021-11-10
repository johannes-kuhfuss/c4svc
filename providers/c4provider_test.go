package providers

import (
	"net/http"
	"os"
	"testing"

	"github.com/johannes-kuhfuss/c4svc/config"
	"github.com/johannes-kuhfuss/c4svc/utils/logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const (
	testUrlGood = "https://mediajku.blob.core.windows.net/media-test/TestBild.tif"
	testUrlBad  = "https://mediajku.blob.core.windows.net/media-test/noexist.tif"
)

func initConfig() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Error("Could not open env file", err)
	}
	config.StorageAccountName = os.Getenv("STORAGE_ACCOUNT_NAME")
	config.StorageAccountKey = os.Getenv("STORAGE_ACCOUNT_KEY")
}

func TestProcessFileNoAccessCred(t *testing.T) {
	c4Id, dstUrl, err := C4Provider.ProcessFile("", false)
	assert.NotNil(t, err)
	assert.Nil(t, c4Id)
	assert.Nil(t, dstUrl)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "No storage account access credentials", err.Message())
}

func TestProcessFileEmptyUrl(t *testing.T) {
	config.StorageAccountName = "dummy"
	config.StorageAccountKey = "dummy"
	c4Id, dstUrl, err := C4Provider.ProcessFile("", false)
	assert.Nil(t, c4Id)
	assert.Nil(t, dstUrl)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Cannot parse source URL", err.Message())
}

func TestProcessFileUrlParseError(t *testing.T) {
	config.StorageAccountName = "dummy"
	config.StorageAccountKey = "dummy"
	dummyUrl := "abcdefg"
	c4Id, dstUrl, err := C4Provider.ProcessFile(dummyUrl, false)
	assert.Nil(t, c4Id)
	assert.Nil(t, dstUrl)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Cannot parse source URL", err.Message())
}

func TestProcessFileWrongCredentials(t *testing.T) {
	config.StorageAccountName = "dummy"
	config.StorageAccountKey = "dummy"
	c4Id, dstUrl, err := C4Provider.ProcessFile(testUrlGood, false)
	assert.Nil(t, c4Id)
	assert.Nil(t, dstUrl)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Cannot access storage account - wrong credentials", err.Message())
}

func TestProcessFileFileNotFoundError(t *testing.T) {
	initConfig()
	c4Id, dstUrl, err := C4Provider.ProcessFile(testUrlBad, false)
	assert.Nil(t, c4Id)
	assert.Nil(t, dstUrl)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Cannot access file on storage account", err.Message())
}

func TestProcessFileNoErrorNoRename(t *testing.T) {
	initConfig()
	c4Id, dstUrl, err := C4Provider.ProcessFile(testUrlGood, false)
	assert.NotNil(t, c4Id)
	assert.Nil(t, dstUrl)
	assert.Nil(t, err)
	assert.EqualValues(t, "c42FTpMRKrEEL6sgwRVRfxbzYDsYZe4VgsNVC7D6Jkqz8ABjsSAybKLYwPLGSJexGkJ9qt3aR8sMAjZ8fhKd7GfQsB", *c4Id)
}

/*
func TestProcessFileNoErrorRename(t *testing.T) {
	initConfig()
	c4Id, dstUrl, err := C4Provider.ProcessFile(testUrlGood, true)
	assert.NotNil(t, c4Id)
	assert.NotNil(t, dstUrl)
	assert.Nil(t, err)
	assert.EqualValues(t, "c42FTpMRKrEEL6sgwRVRfxbzYDsYZe4VgsNVC7D6Jkqz8ABjsSAybKLYwPLGSJexGkJ9qt3aR8sMAjZ8fhKd7GfQsB", *c4Id)
	assert.EqualValues(t, "https://mediajku.blob.core.windows.net/media/c42FTpMRKrEEL6sgwRVRfxbzYDsYZe4VgsNVC7D6Jkqz8ABjsSAybKLYwPLGSJexGkJ9qt3aR8sMAjZ8fhKd7GfQsB.tif", *dstUrl)
}
*/
