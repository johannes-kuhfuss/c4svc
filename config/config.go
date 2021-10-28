package config

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/c4svc/utils/logger"
)

var (
	ginMode            = "debug" // release, debug, test
	ShutDown           = false
	NoJobWaitTime      = (time.Second * 10)
	DeleteFinishedAge  = (time.Hour * 1)
	DeleteFailedAge    = (time.Hour * 2)
	WarnStatusAge      = (time.Hour * 5)
	CleanupWaitTime    = (time.Hour * 1)
	StorageAccountName = ""
	StorageAccountKey  = ""
)

func init() {
	logger.Info("Initalizing configuration")
	osGinMode := os.Getenv("GIN_MODE")
	if osGinMode == gin.ReleaseMode || osGinMode == gin.DebugMode || osGinMode == gin.TestMode {
		ginMode = osGinMode
	}
	StorageAccountName = os.Getenv("STORAGE_ACCOUNT_NAME")
	StorageAccountKey = os.Getenv("STORAGE_ACCOUNT_KEY")
	logger.Info("Done initalizing configuration")
}

func GinMode() string {
	return ginMode
}
