package domain

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
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
	Get(string) (*Job, api_error.ApiErr)
	Save(Job, bool) (*Job, api_error.ApiErr)
	Delete(string) api_error.ApiErr
	GetNext() (*Job, api_error.ApiErr)
	ChangeStatus(string, string) api_error.ApiErr
	CleanJobs(time.Duration, time.Duration) (int, api_error.ApiErr)
	SetC4Id(string, string) api_error.ApiErr
	SetDstUrl(string, string) api_error.ApiErr
	SetErrMsg(string, string) api_error.ApiErr
	GetAll() (*Jobs, api_error.ApiErr)
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

func getJob(jobId string) (*Job, api_error.ApiErr) {
	jobs.mu.Lock()
	defer jobs.mu.Unlock()
	if job := jobs.list[jobId]; job != nil {
		return job, nil
	}
	err := api_error.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return nil, err
}

func (jd *jobDao) Get(jobId string) (*Job, api_error.ApiErr) {
	getJob, err := getJob(jobId)
	if err != nil {
		return nil, err
	}
	return getJob, nil
}

func (jd *jobDao) Save(newJob Job, overwrite bool) (*Job, api_error.ApiErr) {
	_, err := getJob(newJob.Id)
	if err == nil && !overwrite {
		err := api_error.NewBadRequestError(fmt.Sprintf("job with Id %v already exists", newJob.Id))
		return nil, err
	}
	addJob(newJob)
	return &newJob, nil
}

func (jd *jobDao) Delete(jobId string) api_error.ApiErr {
	delJob, _ := getJob(jobId)
	if delJob != nil {
		removeJob(*delJob)
		return nil
	}
	err := api_error.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	return err
}

func (jd *jobDao) GetNext() (*Job, api_error.ApiErr) {
	nextJobId := ""
	nextJobDate := date.GetNowUtc()
	if len(jobs.list) == 0 {
		err := api_error.NewNotFoundError("no jobs in list")
		return nil, err
	}
	for _, v := range jobs.list {
		if v.Status == JobStatusCreated {
			curJobDate, _ := time.Parse(date.ApiDateLayout, v.CreatedAt)
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

func (jd *jobDao) ChangeStatus(jobId string, newStatus string) api_error.ApiErr {
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
		retErr := api_error.NewBadRequestError("invalid status value")
		return retErr
	}
	getJob.ModifiedAt = date.GetNowUtcString()
	_, saveErr := JobDao.Save(*getJob, true)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (jd *jobDao) CleanJobs(finishedTime time.Duration, failedTime time.Duration) (int, api_error.ApiErr) {
	delJobCounter := 0
	if len(jobs.list) == 0 {
		err := api_error.NewNotFoundError("no jobs in list")
		return 0, err
	}
	for _, v := range jobs.list {
		now := date.GetNowUtc()
		modDate, err := time.Parse(date.ApiDateLayout, v.ModifiedAt)
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

func (jd *jobDao) SetC4Id(jobId string, c4Id string) api_error.ApiErr {
	getJob, err := getJob(jobId)
	if err != nil {
		return err
	}
	if strings.TrimSpace(c4Id) == "" {
		return api_error.NewBadRequestError("invalid C4 Id")
	}
	getJob.FileC4Id = c4Id
	getJob.ModifiedAt = date.GetNowUtcString()
	_, saveErr := JobDao.Save(*getJob, true)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (jd *jobDao) SetDstUrl(jobId string, dstUrl string) api_error.ApiErr {
	getJob, err := getJob(jobId)
	if err != nil {
		return err
	}
	if strings.TrimSpace(dstUrl) == "" {
		return api_error.NewBadRequestError("invalid destination URL")
	}
	getJob.DstUrl = dstUrl
	getJob.ModifiedAt = date.GetNowUtcString()
	_, saveErr := JobDao.Save(*getJob, true)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (jd *jobDao) SetErrMsg(jobId string, errMsg string) api_error.ApiErr {
	getJob, err := getJob(jobId)
	if err != nil {
		return err
	}
	getJob.ErrorMsg = errMsg
	getJob.ModifiedAt = date.GetNowUtcString()
	_, saveErr := JobDao.Save(*getJob, true)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (jd *jobDao) GetAll() (*Jobs, api_error.ApiErr) {
	if len(jobs.list) == 0 {
		return nil, api_error.NewNotFoundError("no jobs in list")
	}
	var returnJobs Jobs
	jobs.mu.Lock()
	defer jobs.mu.Unlock()
	for job := range jobs.list {
		returnJobs = append(returnJobs, *jobs.list[job])
	}
	return &returnJobs, nil
}
