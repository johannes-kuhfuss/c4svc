package config

import (
	"os"

	"github.com/gin-gonic/gin"
)

var (
	ginMode = "debug" // release, debug, test
)

func init() {
	osGinMode := os.Getenv("GIN_MODE")
	if osGinMode == gin.ReleaseMode || osGinMode == gin.DebugMode || osGinMode == gin.TestMode {
		ginMode = osGinMode
	}
}

func GinMode() string {
	return ginMode
}
