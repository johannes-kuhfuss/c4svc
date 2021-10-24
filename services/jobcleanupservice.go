package services

import (
	"time"

	"github.com/johannes-kuhfuss/c4/config"
)

var (
	JobCleanupService jobCleanupServiceInterface = &jobCleanupService{}
)

type jobCleanupService struct{}

type jobCleanupServiceInterface interface {
	Cleanup()
}

func (jc *jobCleanupService) Cleanup() {
	for !config.ShutDown {
		time.Sleep(time.Hour * 1)
	}
}
