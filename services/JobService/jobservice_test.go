package services

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	domain "github.com/johannes-kuhfuss/c4/domain/job"
	"github.com/johannes-kuhfuss/c4/utils/date_utils"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	//jobDaoMock      jobsDaoMock
	getJobFunction    func(jobId string) (*domain.Job, rest_errors.RestErr)
	saveJobFunction   func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr)
	deleteJobFunction func(jobId string) rest_errors.RestErr
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
	return deleteJobFunction(jobId)
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

func TestCreateJobInvalidJobType(t *testing.T) {
	newJob := domain.Job{
		Type: "Does not exist",
	}
	createJob, err := JobService.Create(newJob)
	assert.Nil(t, createJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "invalid job type", err.Message())
}

func TestCreateJobInvalidSrcUrl(t *testing.T) {
	newJob := domain.Job{
		Type:   "Create",
		SrcUrl: "",
	}
	createJob, err := JobService.Create(newJob)
	assert.Nil(t, createJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "invalid source Url", err.Message())
}

func TestCreateJobInvalidDstUrl(t *testing.T) {
	newJob := domain.Job{
		Type:   "CreateAndRename",
		SrcUrl: "http://server/path/file.ext",
		DstUrl: "",
	}
	createJob, err := JobService.Create(newJob)
	assert.Nil(t, createJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "invalid destination Url", err.Message())
}

func TestCreateJobNameGivenNoDstUrlNoError(t *testing.T) {
	saveJobFunction = func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr) {
		return &newJob, nil
	}
	newJob := domain.Job{
		Type:   "Create",
		SrcUrl: "http://server/path/file.ext",
		Name:   "myJob",
	}
	createJob, err := JobService.Create(newJob)
	assert.NotNil(t, createJob)
	assert.Nil(t, err)
	_, parseErr := ksuid.Parse(createJob.Id)
	assert.True(t, parseErr == nil)
	assert.EqualValues(t, "myJob", createJob.Name)
	_, parseErr = time.Parse(date_utils.ApiDateLayout, createJob.CreatedAt)
	assert.True(t, parseErr == nil)
	assert.EqualValues(t, newJob.SrcUrl, createJob.SrcUrl)
	assert.EqualValues(t, newJob.Type, createJob.Type)
	assert.EqualValues(t, "", createJob.DstUrl)
}

func TestCreateJobNoNameGivenWithDstUrlNoError(t *testing.T) {
	saveJobFunction = func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr) {
		return &newJob, nil
	}
	newJob := domain.Job{
		Type:   "CreateAndRename",
		SrcUrl: "http://server/path/file.ext",
		DstUrl: "http://server2/path2/file.ext",
	}
	createJob, err := JobService.Create(newJob)
	assert.NotNil(t, createJob)
	assert.Nil(t, err)
	assert.Contains(t, createJob.Name, "Job @ ")
	assert.EqualValues(t, newJob.Type, createJob.Type)
	assert.EqualValues(t, newJob.DstUrl, createJob.DstUrl)
}

func TestCreateJobSaveError(t *testing.T) {
	saveJobFunction = func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr) {
		return nil, rest_errors.NewBadRequestError("could not save job")
	}
	newJob := domain.Job{
		Type:   "Create",
		SrcUrl: "http://server/path/file.ext",
		Name:   "myJob",
	}
	createJob, err := JobService.Create(newJob)
	assert.Nil(t, createJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "could not save job", err.Message())
}

func TestDeleteJobNotFound(t *testing.T) {
	deleteJobFunction = func(jobId string) rest_errors.RestErr {
		return rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	}
	deleteErr := JobService.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, deleteErr)
	assert.EqualValues(t, http.StatusNotFound, deleteErr.StatusCode())
	assert.EqualValues(t, "job with Id 1zXgBZNnBG1msmF1ARQK9ZphbbO does not exist", deleteErr.Message())
}

func TestDeleteJobNoError(t *testing.T) {
	deleteJobFunction = func(jobId string) rest_errors.RestErr {
		return nil
	}
	deleteErr := JobService.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.Nil(t, deleteErr)
}
