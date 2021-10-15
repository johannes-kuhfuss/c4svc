package domain

import (
	"fmt"

	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

var (
	jobs = map[int64]*C4job{
		1: {Id: 1,
			Name:      "Job 1",
			CreatedAt: "2021-10-15T15:00:00Z",
			CreatedBy: "user A",
			SrcUrl:    "http://server/path1/file1.ext",
			Type:      "JobTypeCreate",
			Status:    "Running",
		},
	}

	C4jobDao c4jobDaoInterface = &c4jobDao{}
)

type c4jobDaoInterface interface {
	GetJob(int64) (*C4job, rest_errors.RestErr)
	SaveJob(job C4job) (*C4job, rest_errors.RestErr)
}

type c4jobDao struct{}

func (job *c4jobDao) GetJob(jobId int64) (*C4job, rest_errors.RestErr) {
	if job := jobs[jobId]; job != nil {
		return job, nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job %v does not exist", jobId))
	return nil, err
}

func (job *c4jobDao) SaveJob(newJob C4job) (*C4job, rest_errors.RestErr) {
	if _, found := jobs[newJob.Id]; found {
		err := rest_errors.NewBadRequestError(fmt.Sprintf("job with Id %v already exists", newJob.Id))
		return nil, err
	}
	jobs[newJob.Id] = &newJob
	return &newJob, nil
}
