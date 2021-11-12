package config

import (
	"fmt"
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
	ListenAddr         = ""
)

func init() {
	logger.Info("Initalizing configuration")
	/*
		err := godotenv.Load(".env")
		if err != nil {
			logger.Error("Could not open env file", err)
		}
	*/
	osGinMode := os.Getenv("GIN_MODE")
	if osGinMode == gin.ReleaseMode || osGinMode == gin.DebugMode || osGinMode == gin.TestMode {
		ginMode = osGinMode
	}
	logger.Debug(fmt.Sprintf("Gin-Gonic Mode: %v\n", ginMode))
	StorageAccountName = os.Getenv("STORAGE_ACCOUNT_NAME")
	StorageAccountKey = os.Getenv("STORAGE_ACCOUNT_KEY")
	logger.Debug(fmt.Sprintf("Storage Account Name: %v\n", StorageAccountName))
	ListenAddr = os.Getenv("LISTEN_ADDR")
	if len(ListenAddr) == 0 {
		ListenAddr = ":8080"
	}
	logger.Info("Done initalizing configuration")
}

func GinMode() string {
	return ginMode
}
