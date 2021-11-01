package services

import (
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/c4svc/src/config"
	"github.com/johannes-kuhfuss/c4svc/src/providers"
	"github.com/johannes-kuhfuss/c4svc/src/utils/logger"
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
			rename := curJob.Type == "CreateAndRename"
			c4Id, dstUrl, err := providers.C4Provider.ProcessFile(curJob.SrcUrl, rename)
			if err != nil {
				logger.Error("could process file", err)
				err = JobService.SetErrMsg(curJob.Id, fmt.Sprintf("Could not process file: %s", err.Message()))
				if err != nil {
					logger.Error("could not set error message", err)
				}
				err = JobService.ChangeStatus(curJob.Id, "Failed")
				if err != nil {
					logger.Error("could not change job status", err)
				}
			} else {
				err = JobService.SetC4Id(curJob.Id, *c4Id)
				if err != nil {
					logger.Error("could not set C4 Id", err)
				}
				if rename {
					err = JobService.SetDstUrl(curJob.Id, *dstUrl)
					if err != nil {
						logger.Error("could not set destination URL", err)
					}
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
