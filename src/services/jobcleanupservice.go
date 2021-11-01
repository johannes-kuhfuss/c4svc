package services

import (
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/c4svc/src/config"
	"github.com/johannes-kuhfuss/c4svc/src/domain"
	"github.com/johannes-kuhfuss/c4svc/src/utils/logger"
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
		time.Sleep(config.CleanupWaitTime)
		jobsCleaned, err := domain.JobDao.CleanJobs(config.DeleteFinishedAge, config.DeleteFailedAge)
		if err != nil {
			logger.Info(err.Message())
		} else {
			logger.Info(fmt.Sprintf("Removed %d jobs in state Finished or Failed", jobsCleaned))
		}
	}
}
