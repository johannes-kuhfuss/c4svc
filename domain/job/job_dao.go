package domain

import (
	"fmt"

	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

var (
	jobs = map[string]*Job{
		"1zXgBZNnBG1msmF1ARQK9ZphbbO": {
			Id:         "1zXgBZNnBG1msmF1ARQK9ZphbbO",
			Name:       "Job 1",
			CreatedAt:  "2021-10-15T15:00:00Z",
			CreatedBy:  "user A",
			ModifiedAt: "",
			ModifiedBy: "",
			SrcUrl:     "http://server/path1/file1.ext",
			DstUrl:     "",
			Type:       "JobTypeCreate",
			Status:     "Running",
			FileC4Id:   "abcdefg",
		},
	}

	JobDao jobDaoInterface = &jobDao{}
)

type jobDaoInterface interface {
	Get(string) (*Job, rest_errors.RestErr)
	Save(Job, bool) (*Job, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
}

type jobDao struct{}

func (job *jobDao) Get(jobId string) (*Job, rest_errors.RestErr) {
	if job := jobs[jobId]; job != nil {
		return job, nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return nil, err
}

func (job *jobDao) Save(newJob Job, overwrite bool) (*Job, rest_errors.RestErr) {
	_, found := jobs[newJob.Id]
	if found && !overwrite {
		err := rest_errors.NewBadRequestError(fmt.Sprintf("job with Id %v already exists", newJob.Id))
		return nil, err
	}
	jobs[newJob.Id] = &newJob
	return &newJob, nil
}

func (job *jobDao) Delete(jobId string) rest_errors.RestErr {
	if job := jobs[jobId]; job != nil {
		delete(jobs, jobId)
		return nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return err
}
