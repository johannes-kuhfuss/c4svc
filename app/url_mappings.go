package app

import (
	"github.com/johannes-kuhfuss/c4/controllers/job"
	"github.com/johannes-kuhfuss/c4/controllers/ping"
)

func mapUrls() {
	router.GET("/ping", ping.Pong)
	router.POST("/job", job.CreateJob)
	router.GET("/job/:jobid", job.GetJob)
}
