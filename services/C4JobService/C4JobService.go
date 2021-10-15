package services

import (
	"fmt"
	"strings"

	domain "github.com/johannes-kuhfuss/c4/domain/c4job"
	"github.com/johannes-kuhfuss/c4/utils/date_utils"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

type c4JobService struct{}

type c4JobServiceInterface interface {
	CreateC4Job(domain.C4job) (*domain.C4job, rest_errors.RestErr)
	GetC4Job(int64) (*domain.C4job, rest_errors.RestErr)
}

var (
	C4JobService c4JobServiceInterface
)

func Init() {
	C4JobService = &c4JobService{}
}

func (c4 *c4JobService) CreateC4Job(inputJob domain.C4job) (*domain.C4job, rest_errors.RestErr) {
	if err := inputJob.Validate(); err != nil {
		return nil, err
	}
	request := domain.C4job{}
	request.Id = 2
	if strings.TrimSpace(inputJob.Name) != "" {
		request.Name = inputJob.Name
	} else {
		request.Name = fmt.Sprintf("C4 Job @ %s", date_utils.GetNowUtcString())
	}
	request.CreatedAt = date_utils.GetNowUtcString()
	request.CreatedBy = "user-im"
	request.SrcUrl = inputJob.SrcUrl
	if inputJob.Type == domain.JobTypeCreateAndRename {
		request.DstUrl = inputJob.DstUrl
	} else {
		request.DstUrl = ""
	}
	request.Type = inputJob.Type
	request.Status = domain.JobStatusCreated
	savedJob, err := domain.C4jobDao.SaveJob(request)
	if err != nil {
		return nil, err
	}
	return savedJob, nil
}

func (c4 *c4JobService) GetC4Job(id int64) (*domain.C4job, rest_errors.RestErr) {
	return nil, nil
}
