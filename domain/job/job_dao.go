package domain

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/johannes-kuhfuss/c4/utils/date_utils"
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
	GetNext() (*Job, rest_errors.RestErr)
	ChangeStatus(string, string) rest_errors.RestErr
	CleanJobs(time.Duration, time.Duration) (int, rest_errors.RestErr)
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

func (jd *jobDao) Get(jobId string) (*Job, rest_errors.RestErr) {
	getJob, err := getJob(jobId)
	if err != nil {
		return nil, err
	}
	return getJob, nil
}

func (jd *jobDao) Save(newJob Job, overwrite bool) (*Job, rest_errors.RestErr) {
	_, err := getJob(newJob.Id)
	if err == nil && !overwrite {
		err := rest_errors.NewBadRequestError(fmt.Sprintf("job with Id %v already exists", newJob.Id))
		return nil, err
	}
	addJob(newJob)
	return &newJob, nil
}

func (jd *jobDao) Delete(jobId string) rest_errors.RestErr {
	delJob, _ := getJob(jobId)
	if delJob != nil {
		removeJob(*delJob)
		return nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return err
}

func (jd *jobDao) GetNext() (*Job, rest_errors.RestErr) {
	nextJobId := ""
	nextJobDate := date_utils.GetNowUtc()
	if len(jobs.list) == 0 {
		err := rest_errors.NewNotFoundError("no jobs in list")
		return nil, err
	}
	for _, v := range jobs.list {
		if v.Status == JobStatusCreated {
			curJobDate, _ := time.Parse(date_utils.ApiDateLayout, v.CreatedAt)
			if curJobDate.Before(nextJobDate) {
				nextJobDate = curJobDate
				nextJobId = v.Id
			}
		}
	}
	getJob, err := getJob(nextJobId)
	if err != nil {
		return nil, err
	}
	return getJob, nil
}

func (jd *jobDao) ChangeStatus(jobId string, newStatus string) rest_errors.RestErr {
	getJob, err := getJob(jobId)
	if err != nil {
		return err
	}
	newStatus = strings.ToLower(newStatus)
	if strings.ToLower(string(getJob.Status)) == newStatus {
		return nil
	}
	switch newStatus {
	case "created":
		getJob.Status = JobStatusCreated
	case "running":
		getJob.Status = JobStatusRunning
	case "failed":
		getJob.Status = JobStatusFailed
	case "finished":
		getJob.Status = JobStatusFinished
	default:
		retErr := rest_errors.NewBadRequestError("invalid status value")
		return retErr
	}
	getJob.ModifiedAt = date_utils.GetNowUtcString()
	_, saveErr := JobDao.Save(*getJob, true)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (jd *jobDao) CleanJobs(finishedTime time.Duration, failedTime time.Duration) (int, rest_errors.RestErr) {
	delJobCounter := 0
	if len(jobs.list) == 0 {
		err := rest_errors.NewNotFoundError("no jobs in list")
		return 0, err
	}
	for _, v := range jobs.list {
		now := date_utils.GetNowUtc()
		modDate, err := time.Parse(date_utils.ApiDateLayout, v.ModifiedAt)
		if err != nil {
			continue
		}
		if (v.Status == JobStatusFailed && modDate.Add(failedTime).Before(now)) || (v.Status == JobStatusFinished && modDate.Add(finishedTime).Before(now)) {
			delJobCounter++
			removeJob(*v)
		}
	}
	return delJobCounter, nil
}
