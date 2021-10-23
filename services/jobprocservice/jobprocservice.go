package services

import (
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/c4/config"
	services "github.com/johannes-kuhfuss/c4/services/jobservice"
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
			// to do: process job
			time.Sleep(time.Second * time.Duration(config.JobCycleTime))
			fmt.Println("found job")
		}
		_ = curJob
		time.Sleep(time.Second * time.Duration(config.JobCycleTime))
	}
}
