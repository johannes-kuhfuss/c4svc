package services

import (
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/c4/config"
	"github.com/johannes-kuhfuss/c4/providers"
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
		curJob, err := JobService.GetNext()
		if err == nil {
			logger.Debug(fmt.Sprintf("Found job with Id %v to process", curJob.Id))
			c4Id, err := providers.C4Provider.ProcessFile(curJob.SrcUrl)
			if err != nil {
				logger.Error("could process file", err)
				err = JobService.ChangeStatus(curJob.Id, "Failed")
				if err != nil {
					logger.Error("could not change job status", err)
				}
			} else {
				fmt.Println(c4Id)
				// to do write c4 back to job
				err = JobService.ChangeStatus(curJob.Id, "Finished")
				if err != nil {
					logger.Error("could not change job status", err)
				}
			}

			logger.Debug(fmt.Sprintf("Done processing job with Id %v", curJob.Id))
		} else {
			logger.Debug("no job found. Sleeping...")
			time.Sleep(time.Second * time.Duration(config.NoJobWaitTime))
		}

	}
}
