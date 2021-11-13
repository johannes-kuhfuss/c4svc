package app

import (
	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/c4svc/config"
	"github.com/johannes-kuhfuss/c4svc/services"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	router *gin.Engine
)

func init() {
	logger.Debug("Initializing router")
	gin.SetMode(config.GinMode())
	gin.DefaultWriter = logger.GetLogger()
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	logger.Debug("Done initializing router")
}

func StartApp() {
	logger.Info("Starting application")
	mapUrls()

	logger.Info("Starting job processor")
	go services.JobProcService.Process()
	logger.Info("Starting job cleanup")
	go services.JobCleanupService.Cleanup()

	if err := router.Run(config.ListenAddr); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}

	logger.Info("Application ended")
}
