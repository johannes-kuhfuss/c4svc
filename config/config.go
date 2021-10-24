package config

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/c4/utils/logger"
)

var (
	ginMode           = "debug" // release, debug, test
	ShutDown          = false
	NoJobWaitTime     = (time.Second * 10)
	DeleteFinishedAge = (time.Hour * 1)
	DeleteFailedAge   = (time.Hour * 2)
	WarnStatusAge     = (time.Hour * 5)
)

func init() {
	logger.Info("Initalizing configuration")
	osGinMode := os.Getenv("GIN_MODE")
	if osGinMode == gin.ReleaseMode || osGinMode == gin.DebugMode || osGinMode == gin.TestMode {
		ginMode = osGinMode
	}
	logger.Info("Done initalizing configuration")
}

func GinMode() string {
	return ginMode
}
