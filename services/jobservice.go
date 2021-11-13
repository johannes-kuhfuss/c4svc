package services

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/c4svc/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/segmentio/ksuid"
)

var (
	JobService jobServiceInterface = &jobService{}
)

type jobService struct{}

type jobServiceInterface interface {
	Create(domain.Job) (*domain.Job, api_error.ApiErr)
	Get(string) (*domain.Job, api_error.ApiErr)
	Delete(string) api_error.ApiErr
	Update(string, domain.Job, bool) (*domain.Job, api_error.ApiErr)
	GetNext() (*domain.Job, api_error.ApiErr)
	ChangeStatus(string, string) api_error.ApiErr
	SetC4Id(string, string) api_error.ApiErr
	SetDstUrl(string, string) api_error.ApiErr
	SetErrMsg(string, string) api_error.ApiErr
	GetAll() (*domain.Jobs, api_error.ApiErr)
}

func (j *jobService) Create(inputJob domain.Job) (*domain.Job, api_error.ApiErr) {
	if err := inputJob.Validate(); err != nil {
		return nil, err
	}
	request := domain.Job{}
	request.Id = ksuid.New().String()
	if strings.TrimSpace(inputJob.Name) != "" {
		request.Name = inputJob.Name
	} else {
		request.Name = fmt.Sprintf("Job @ %s", date.GetNowUtcString())
	}
	request.CreatedAt = date.GetNowUtcString()
	request.SrcUrl = inputJob.SrcUrl
	request.DstUrl = ""
	request.Type = inputJob.Type
	request.Status = domain.JobStatusCreated
	savedJob, err := domain.JobDao.Save(request, false)
	if err != nil {
		return nil, err
	}
	return savedJob, nil
}

func (j *jobService) Get(jobId string) (*domain.Job, api_error.ApiErr) {
	job, err := domain.JobDao.Get(jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *jobService) Delete(jobId string) api_error.ApiErr {
	job, err := domain.JobDao.Get(jobId)
	if err != nil {
		return err
	}
	if job.Status == domain.JobStatusRunning {
		statusErr := api_error.NewProcessingConflictError("Cannot delete job in status running")
		return statusErr
	}
	deleteErr := domain.JobDao.Delete(jobId)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (j *jobService) Update(jobId string, inputJob domain.Job, partial bool) (*domain.Job, api_error.ApiErr) {
	job, err := domain.JobDao.Get(jobId)
	if err != nil {
		return nil, err
	}
	if job.Status != domain.JobStatusCreated {
		statusErr := api_error.NewProcessingConflictError("Cannot modify job in status other than created")
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
	request.ModifiedAt = date.GetNowUtcString()
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

func (j *jobService) GetNext() (*domain.Job, api_error.ApiErr) {
	job, err := domain.JobDao.GetNext()
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *jobService) ChangeStatus(jobId string, newStatus string) api_error.ApiErr {
	err := domain.JobDao.ChangeStatus(jobId, newStatus)
	if err != nil {
		return err
	}
	return nil
}

func (j *jobService) SetC4Id(jobId string, c4Id string) api_error.ApiErr {
	err := domain.JobDao.SetC4Id(jobId, c4Id)
	if err != nil {
		return err
	}
	return nil
}

func (j *jobService) SetDstUrl(jobId string, dstUrl string) api_error.ApiErr {
	err := domain.JobDao.SetDstUrl(jobId, dstUrl)
	if err != nil {
		return err
	}
	return nil
}

func (j *jobService) SetErrMsg(jobId string, errMsg string) api_error.ApiErr {
	err := domain.JobDao.SetErrMsg(jobId, errMsg)
	if err != nil {
		return err
	}
	return nil
}

func (j *jobService) GetAll() (*domain.Jobs, api_error.ApiErr) {
	jobs, err := domain.JobDao.GetAll()
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
