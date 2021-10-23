package domain

import (
	"fmt"
	"sync"

	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

var (
	jobs = jobList{
		list: make(map[string]*Job),
		mu:   sync.Mutex{},
	}
	JobDao jobDaoInterface = &jobDao{}
)

type jobList struct {
	list map[string]*Job
	mu   sync.Mutex
}

type jobDaoInterface interface {
	Get(string) (*Job, rest_errors.RestErr)
	Save(Job, bool) (*Job, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
}

type jobDao struct{}

func addJob(newJob Job) {
	jobs.mu.Lock()
	defer jobs.mu.Unlock()
	jobs.list[newJob.Id] = &newJob
}

func removeJob(delJob Job) {
	jobs.mu.Lock()
	defer jobs.mu.Unlock()
	delete(jobs.list, delJob.Id)
}

func getJob(jobId string) (*Job, rest_errors.RestErr) {
	jobs.mu.Lock()
	defer jobs.mu.Unlock()
	if job := jobs.list[jobId]; job != nil {
		return job, nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return nil, err
}

func (job *jobDao) Get(jobId string) (*Job, rest_errors.RestErr) {
	getJob, err := getJob(jobId)
	if err != nil {
		return nil, err
	}
	return getJob, nil
}

func (job *jobDao) Save(newJob Job, overwrite bool) (*Job, rest_errors.RestErr) {
	_, err := getJob(newJob.Id)
	if err == nil && !overwrite {
		err := rest_errors.NewBadRequestError(fmt.Sprintf("job with Id %v already exists", newJob.Id))
		return nil, err
	}
	addJob(newJob)
	return &newJob, nil
}

func (job *jobDao) Delete(jobId string) rest_errors.RestErr {
	delJob, _ := getJob(jobId)
	if delJob != nil {
		removeJob(*delJob)
		return nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return err
}
