package app

import (
	"github.com/johannes-kuhfuss/c4/controllers/c4job"
	"github.com/johannes-kuhfuss/c4/controllers/ping"
)

func mapUrls() {
	router.GET("/ping", ping.Pong)
	router.POST("/c4job", c4job.CreateC4Job)
	router.GET("/c4job/:jobid", c4job.GetC4Job)
}
