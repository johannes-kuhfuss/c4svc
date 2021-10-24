package app

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/c4/config"
	srv2 "github.com/johannes-kuhfuss/c4/services/jobcleanupservice"
	srv1 "github.com/johannes-kuhfuss/c4/services/jobprocservice"
	"github.com/johannes-kuhfuss/c4/utils/logger"
)

var (
	router *gin.Engine
)

func init() {
	logger.Debug("Initializing router")
	gin.SetMode(config.GinMode())
	router = gin.New()
	router.Use(ginzap.Ginzap(logger.GetLog(), time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger.GetLog(), true))
	logger.Debug("Done initializing router")
}

func StartApp() {
	logger.Info("Starting application")
	mapUrls()

	logger.Info("Starting job processor")
	go srv1.JobProcService.Process()
	logger.Info("Starting job cleanup")
	go srv2.JobCleanupService.Cleanup()

	if err := router.Run(":8080"); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}

	logger.Info("Application ended")
}
