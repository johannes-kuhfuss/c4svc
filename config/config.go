package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/c4/utils/logger"
)

var (
	ginMode = "debug" // release, debug, test
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
