package services

import (
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/c4svc/config"
	"github.com/johannes-kuhfuss/c4svc/providers"
	"github.com/johannes-kuhfuss/c4svc/utils/logger"
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
		curJob, err := JobService.GetNext()
		if err == nil {
			logger.Info(fmt.Sprintf("Found job with Id %v to process", curJob.Id))
			c4Id, err := providers.C4Provider.ProcessFile(curJob.SrcUrl)
			if err != nil {
				logger.Error("could process file", err)
				err = JobService.ChangeStatus(curJob.Id, "Failed")
				if err != nil {
					logger.Error("could not change job status", err)
				}
			} else {
				err = JobService.SetC4Id(curJob.Id, *c4Id)
				if err != nil {
					logger.Error("could not set C4 Id", err)
				}
				err = JobService.ChangeStatus(curJob.Id, "Finished")
				if err != nil {
					logger.Error("could not change job status", err)
				}
			}

			logger.Info(fmt.Sprintf("Done processing job with Id %v", curJob.Id))
		} else {
			logger.Debug("no job found. Sleeping...")
			time.Sleep(time.Second * time.Duration(config.NoJobWaitTime))
		}

	}
}
