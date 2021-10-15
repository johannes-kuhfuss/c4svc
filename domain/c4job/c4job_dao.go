package domain

import (
	"fmt"

	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

var (
	jobs = map[int64]*C4job{
		1: {Id: 1},
	}

	C4jobDao c4jobDaoInterface
)

func init() {
	C4jobDao = &c4jobDao{}
}

type c4jobDaoInterface interface {
	GetJob(int64) (*C4job, *rest_errors.RestErr)
	SaveJob(job C4job) (*C4job, *rest_errors.RestErr)
}

type c4jobDao struct{}

func (job *c4jobDao) GetJob(jobId int64) (*C4job, *rest_errors.RestErr) {
	if job := jobs[jobId]; job != nil {
		return job, nil
	}
	err := rest_errors.NewNotFoundError(fmt.Sprintf("job %v does not exists", jobId))
	return nil, &err
}

func (job *c4jobDao) SaveJob(newJob C4job) (*C4job, *rest_errors.RestErr) {
	return nil, nil
}
