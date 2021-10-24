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
	Create(domain.Job) (*domain.Job, rest_errors.RestErr)
	Get(string) (*domain.Job, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Update(string, domain.Job, bool) (*domain.Job, rest_errors.RestErr)
	GetNext() (*domain.Job, rest_errors.RestErr)
	ChangeStatus(string, string) rest_errors.RestErr
}

func (j *jobService) Create(inputJob domain.Job) (*domain.Job, rest_errors.RestErr) {
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
	savedJob, err := domain.JobDao.Save(request, false)
	if err != nil {
		return nil, err
	}
	return savedJob, nil
}

func (j *jobService) Get(jobId string) (*domain.Job, rest_errors.RestErr) {
	job, err := domain.JobDao.Get(jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *jobService) Delete(jobId string) rest_errors.RestErr {
	job, err := domain.JobDao.Get(jobId)
	if err != nil {
		return err
	}
	if job.Status == domain.JobStatusRunning {
		statusErr := rest_errors.NewProcessingConflictError("Cannot delete job in status running")
		return statusErr
	}
	deleteErr := domain.JobDao.Delete(jobId)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (j *jobService) Update(jobId string, inputJob domain.Job, partial bool) (*domain.Job, rest_errors.RestErr) {
	job, err := domain.JobDao.Get(jobId)
	if err != nil {
		return nil, err
	}
	if job.Status != domain.JobStatusCreated {
		statusErr := rest_errors.NewProcessingConflictError("Cannot modify job in status other than created")
		return nil, statusErr
	}
	if !partial {
		if err := inputJob.Validate(); err != nil {
			return nil, err
		}
	}
	request := domain.Job{}
	request.Id = job.Id
	request.CreatedAt = job.CreatedAt
	request.CreatedBy = job.CreatedBy
	request.ModifiedAt = date_utils.GetNowUtcString()
	request.Status = job.Status
	request.FileC4Id = job.FileC4Id
	if partial && strings.TrimSpace(inputJob.Name) == "" {
		request.Name = job.Name
	} else {
		request.Name = inputJob.Name
	}
	if partial && strings.TrimSpace(inputJob.SrcUrl) == "" {
		request.SrcUrl = job.SrcUrl
	} else {
		request.SrcUrl = inputJob.SrcUrl
	}
	if partial && strings.TrimSpace(inputJob.DstUrl) == "" {
		request.DstUrl = job.DstUrl
	} else {
		request.DstUrl = inputJob.DstUrl
	}
	if partial && (inputJob.Type == "") {
		request.Type = job.Type
	} else {
		request.Type = inputJob.Type
	}

	savedJob, err := domain.JobDao.Save(request, true)
	if err != nil {
		return nil, err
	}
	return savedJob, nil
}

func (j *jobService) GetNext() (*domain.Job, rest_errors.RestErr) {
	job, err := domain.JobDao.GetNext()
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *jobService) ChangeStatus(jobId string, newStatus string) rest_errors.RestErr {
	err := domain.JobDao.ChangeStatus(jobId, newStatus)
	if err != nil {
		return err
	}
	return nil
}
