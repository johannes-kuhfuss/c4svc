package domain

import (
	"fmt"

	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

var (
	jobs                   = map[string]*Job{}
	JobDao jobDaoInterface = &jobDao{}
)

type jobDaoInterface interface {
	Get(string) (*Job, rest_errors.RestErr)
	Save(Job, bool) (*Job, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
}

type jobDao struct{}

func addJob(newJob Job) {
	jobs[newJob.Id] = &newJob
}

func removeJob(delJob Job) {
	delete(jobs, delJob.Id)
}

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
	addJob(newJob)
	return &newJob, nil
}

func (job *jobDao) Delete(jobId string) rest_errors.RestErr {
	if delJob := jobs[jobId]; delJob != nil {
		removeJob(*delJob)
		return nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return err
}
