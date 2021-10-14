package services

import (
	"github.com/johannes-kuhfuss/c4/domain/c4job"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

type c4JobService struct{}

type c4JobServiceInterface interface {
	CreateC4Job(c4job.C4job) (*c4job.C4job, rest_errors.RestErr)
	GetC4Job(int64) (*c4job.C4job, rest_errors.RestErr)
}

var (
	C4JobService c4JobServiceInterface
)

func Init() {
	C4JobService = &c4JobService{}
}

func (c4 c4JobService) CreateC4Job(job c4job.C4job) (*c4job.C4job, rest_errors.RestErr) {
	return nil, nil
}

func (c4 c4JobService) GetC4Job(id int64) (*c4job.C4job, rest_errors.RestErr) {
	return nil, nil
}
