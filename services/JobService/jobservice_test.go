package services

import (
	"net/http"
	"testing"

	domain "github.com/johannes-kuhfuss/c4/domain/job"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
	"github.com/stretchr/testify/assert"
)

var (
	//jobDaoMock      jobsDaoMock
	getJobFunction  func(jobId string) (*domain.Job, rest_errors.RestErr)
	saveJobFunction func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr)
	deleteJobFunc   func(jobId string) rest_errors.RestErr
)

type jobsDaoMock struct{}

func init() {
	domain.JobDao = &jobsDaoMock{}
}

func (m *jobsDaoMock) Get(jobId string) (*domain.Job, rest_errors.RestErr) {
	return getJobFunction(jobId)
}

func (m *jobsDaoMock) Save(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr) {
	return saveJobFunction(newJob, overwrite)
}

func (m *jobsDaoMock) Delete(jobId string) rest_errors.RestErr {
	return deleteJobFunc(jobId)
}

func TestGetJobNotFound(t *testing.T) {
	getJobFunction = func(jobId string) (*domain.Job, rest_errors.RestErr) {
		return nil, rest_errors.NewNotFoundError("job with Id X does not exist")
	}
	user, err := JobService.Get("X")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "job with Id X does not exist", err.Message())
}

func TestGetJobNoError(t *testing.T) {
	getJobFunction = func(jobId string) (*domain.Job, rest_errors.RestErr) {
		return &domain.Job{
			Id:         jobId,
			Name:       "Job 1",
			CreatedAt:  "2021-10-15T15:00:00Z",
			CreatedBy:  "user A",
			ModifiedAt: "",
			ModifiedBy: "",
			SrcUrl:     "http://server/path1/file1.ext",
			DstUrl:     "",
			Type:       "Create",
			Status:     "Running",
			FileC4Id:   "abcdefg",
		}, nil
	}
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	user, err := JobService.Get(id)
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, user.Id, id)
}
