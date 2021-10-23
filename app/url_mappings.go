package app

import (
	"github.com/johannes-kuhfuss/c4/controllers/job"
	"github.com/johannes-kuhfuss/c4/controllers/ping"
	"github.com/johannes-kuhfuss/c4/utils/logger"
)

func mapUrls() {
	logger.Debug("Mapping URLs")

	router.GET("/ping", ping.Pong)
	router.POST("/job", job.Create)
	router.GET("/job/:job_id", job.Get)
	router.DELETE("/job/:job_id", job.Delete)
	router.PUT("/job/:job_id", job.Update)
	router.PATCH("/job/:job_id", job.UpdatePart)

	logger.Debug("Done mapping URLs")
}
