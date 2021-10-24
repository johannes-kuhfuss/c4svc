package services

import (
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/c4/config"
	services "github.com/johannes-kuhfuss/c4/services/jobservice"
	"github.com/johannes-kuhfuss/c4/utils/logger"
)

var (
	JobProcService jobProcServiceInterface = &jobProcService{}
)

type jobProcService struct{}

type jobProcServiceInterface interface {
	Process()
}

func (jp *jobProcService) Process() {
	for !config.ShutDown {
		curJob, err := services.JobService.GetNext()
		if err == nil {
			logger.Debug(fmt.Sprintf("Found job with Id %v to process", curJob.Id))
			// to do: process job
			time.Sleep(time.Second * 5)
			err = services.JobService.ChangeStatus(curJob.Id, "Finished")
			if err != nil {
				logger.Error("could not change job status", err)
			}
			logger.Debug(fmt.Sprintf("Done processing job with Id %v", curJob.Id))
		} else {
			logger.Debug("no job found. Sleeping...")
			time.Sleep(time.Second * time.Duration(config.JobCycleTime))
		}

	}
}
