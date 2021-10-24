package services

import "github.com/johannes-kuhfuss/c4/utils/logger"

var (
	JobCleanupService jobCleanupServiceInterface = &jobCleanupService{}
)

type jobCleanupService struct{}

type jobCleanupServiceInterface interface {
	Cleanup()
}

func (jc *jobCleanupService) Cleanup() {
	logger.Info("Cleaning up")
}
