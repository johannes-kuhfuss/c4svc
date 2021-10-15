package services

import (
	"fmt"
	"strings"

	domain "github.com/johannes-kuhfuss/c4/domain/job"
	"github.com/johannes-kuhfuss/c4/utils/date_utils"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
	"github.com/segmentio/ksuid"
)

var (
	JobService jobServiceInterface = &jobService{}
)

type jobService struct{}

type jobServiceInterface interface {
	CreateJob(domain.Job) (*domain.Job, rest_errors.RestErr)
	GetJob(string) (*domain.Job, rest_errors.RestErr)
}

func (j *jobService) CreateJob(inputJob domain.Job) (*domain.Job, rest_errors.RestErr) {
	if err := inputJob.Validate(); err != nil {
		return nil, err
	}
	request := domain.Job{}
	request.Id = ksuid.New().String()
	if strings.TrimSpace(inputJob.Name) != "" {
		request.Name = inputJob.Name
	} else {
		request.Name = fmt.Sprintf("Job @ %s", date_utils.GetNowUtcString())
	}
	request.CreatedAt = date_utils.GetNowUtcString()
	request.SrcUrl = inputJob.SrcUrl
	if inputJob.Type == domain.JobTypeCreateAndRename {
		request.DstUrl = inputJob.DstUrl
	} else {
		request.DstUrl = ""
	}
	request.Type = inputJob.Type
	request.Status = domain.JobStatusCreated
	savedJob, err := domain.JobDao.SaveJob(request)
	if err != nil {
		return nil, err
	}
	return savedJob, nil
}

func (j *jobService) GetJob(jobId string) (*domain.Job, rest_errors.RestErr) {
	job, err := domain.JobDao.GetJob(jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}
