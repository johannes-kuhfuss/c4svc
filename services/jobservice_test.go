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
	getJobFunction       func(jobId string) (*domain.Job, rest_errors.RestErr)
	saveJobFunction      func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr)
	deleteJobFunction    func(jobId string) rest_errors.RestErr
	getNextJobFunction   func() (*domain.Job, rest_errors.RestErr)
	changeStatusFunction func(jobId string, newStatus string) rest_errors.RestErr
	cleanJobsFunction    func(finishedTime time.Duration, failedTime time.Duration) (int, rest_errors.RestErr)
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

func (m *jobsDaoMock) GetNext() (*domain.Job, rest_errors.RestErr) {
	return getNextJobFunction()
}

func (m *jobsDaoMock) ChangeStatus(jobId string, newStatus string) rest_errors.RestErr {
	return changeStatusFunction(jobId, newStatus)
}

func (m *jobsDaoMock) CleanJobs(finishedTime time.Duration, failedTime time.Duration) (int, rest_errors.RestErr) {
	return cleanJobsFunction(finishedTime, failedTime)
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
			Status:     "Created",
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
		Type: "invalid_Type",
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
	getJobFunction = func(jobId string) (*domain.Job, rest_errors.RestErr) {
		return nil, rest_errors.NewNotFoundError("job with Id 1zXgBZNnBG1msmF1ARQK9ZphbbO does not exist")
	}
	deleteJobFunction = func(jobId string) rest_errors.RestErr {
		return rest_errors.NewNotFoundError(fmt.Sprintf("job with Id %v does not exist", jobId))
	}
	deleteErr := JobService.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, deleteErr)
	assert.EqualValues(t, http.StatusNotFound, deleteErr.StatusCode())
	assert.EqualValues(t, "job with Id 1zXgBZNnBG1msmF1ARQK9ZphbbO does not exist", deleteErr.Message())
}

func TestDeleteJobStatusError(t *testing.T) {
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
	deleteJobFunction = func(jobId string) rest_errors.RestErr {
		return nil
	}
	deleteErr := JobService.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, deleteErr)
	assert.EqualValues(t, http.StatusConflict, deleteErr.StatusCode())
	assert.EqualValues(t, "Cannot delete job in status running", deleteErr.Message())
}

func TestDeleteDeleteError(t *testing.T) {
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
			Status:     "Created",
			FileC4Id:   "abcdefg",
		}, nil
	}
	deleteJobFunction = func(jobId string) rest_errors.RestErr {
		return rest_errors.NewInternalServerError("could not delete job", nil)
	}
	deleteErr := JobService.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, deleteErr)
	assert.EqualValues(t, http.StatusInternalServerError, deleteErr.StatusCode())
	assert.EqualValues(t, "could not delete job", deleteErr.Message())
}

func TestDeleteJobNoError(t *testing.T) {
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
			Status:     "Created",
			FileC4Id:   "abcdefg",
		}, nil
	}
	deleteJobFunction = func(jobId string) rest_errors.RestErr {
		return nil
	}
	deleteErr := JobService.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.Nil(t, deleteErr)
}

func TestUpdateJobNotFound(t *testing.T) {
	getJobFunction = func(jobId string) (*domain.Job, rest_errors.RestErr) {
		return nil, rest_errors.NewNotFoundError("job not found")
	}
	inputJob := domain.Job{}
	updateJob, err := JobService.Update("", inputJob, false)
	assert.Nil(t, updateJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "job not found", err.Message())
}

func TestUpdateJobValidateFailure(t *testing.T) {
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
			Status:     "Created",
			FileC4Id:   "abcdefg",
		}, nil
	}
	inputJob := domain.Job{
		Name:       "",
		CreatedAt:  "2022-11-16T16:01:01Z",
		CreatedBy:  "user B",
		ModifiedAt: "2022-11-16T16:01:01Z",
		ModifiedBy: "user C",
		SrcUrl:     "http://server3/path3/file3.ext",
		DstUrl:     "http://server2/path2/file2.ext",
		Type:       "not_valid",
		Status:     "Created",
		FileC4Id:   "xyz",
	}
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	updateJob, err := JobService.Update(id, inputJob, false)
	assert.Nil(t, updateJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "invalid job type", err.Message())
}

func TestUpdateJobStatusError(t *testing.T) {
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
			Status:     "Failed",
			FileC4Id:   "abcdefg",
		}, nil
	}
	inputJob := domain.Job{
		Name:       "",
		CreatedAt:  "2022-11-16T16:01:01Z",
		CreatedBy:  "user B",
		ModifiedAt: "2022-11-16T16:01:01Z",
		ModifiedBy: "user C",
		SrcUrl:     "http://server3/path3/file3.ext",
		DstUrl:     "http://server2/path2/file2.ext",
		Type:       "not_valid",
		Status:     "Failed",
		FileC4Id:   "xyz",
	}
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	updateJob, err := JobService.Update(id, inputJob, false)
	assert.Nil(t, updateJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusConflict, err.StatusCode())
	assert.EqualValues(t, "Cannot modify job in status other than created", err.Message())
}

func TestUpdateJobFullUpdate(t *testing.T) {
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
			Status:     "Created",
			FileC4Id:   "abcdefg",
		}, nil
	}
	inputJob := domain.Job{
		Name:       "",
		CreatedAt:  "2022-11-16T16:01:01Z",
		CreatedBy:  "user B",
		ModifiedAt: "2022-11-16T16:01:01Z",
		ModifiedBy: "user C",
		SrcUrl:     "http://server3/path3/file3.ext",
		DstUrl:     "http://server2/path2/file2.ext",
		Type:       "CreateAndRename",
		Status:     "Running",
		FileC4Id:   "xyz",
	}
	saveJobFunction = func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr) {
		return &newJob, nil
	}
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	updateJob, err := JobService.Update(id, inputJob, false)
	assert.NotNil(t, updateJob)
	assert.Nil(t, err)
	assert.EqualValues(t, id, updateJob.Id)
	assert.EqualValues(t, inputJob.Name, updateJob.Name)
	assert.EqualValues(t, "2021-10-15T15:00:00Z", updateJob.CreatedAt)
	assert.EqualValues(t, "user A", updateJob.CreatedBy)
	assert.NotEqualValues(t, "", updateJob.ModifiedAt)
	assert.EqualValues(t, "http://server3/path3/file3.ext", updateJob.SrcUrl)
	assert.EqualValues(t, "http://server2/path2/file2.ext", updateJob.DstUrl)
	assert.EqualValues(t, "CreateAndRename", updateJob.Type)
	assert.EqualValues(t, "Created", updateJob.Status)
	assert.EqualValues(t, "abcdefg", updateJob.FileC4Id)
}

func TestUpdateJobPartialUpdate(t *testing.T) {
	getJobFunction = func(jobId string) (*domain.Job, rest_errors.RestErr) {
		return &domain.Job{
			Id:         jobId,
			Name:       "Job 1",
			CreatedAt:  "2021-10-15T15:00:00Z",
			CreatedBy:  "user A",
			ModifiedAt: "",
			ModifiedBy: "",
			SrcUrl:     "http://server/path1/file1.ext",
			DstUrl:     "http://server/path2/file2.ext",
			Type:       "CreateAndRename",
			Status:     "Created",
			FileC4Id:   "abcdefg",
		}, nil
	}
	inputJob := domain.Job{
		Name:   "",
		DstUrl: "",
	}
	saveJobFunction = func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr) {
		return &newJob, nil
	}
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	updateJob, err := JobService.Update(id, inputJob, true)
	assert.NotNil(t, updateJob)
	assert.Nil(t, err)
	assert.EqualValues(t, id, updateJob.Id)
	assert.EqualValues(t, "Job 1", updateJob.Name)
	assert.EqualValues(t, "2021-10-15T15:00:00Z", updateJob.CreatedAt)
	assert.EqualValues(t, "user A", updateJob.CreatedBy)
	assert.NotEqualValues(t, "", updateJob.ModifiedAt)
	assert.EqualValues(t, "http://server/path1/file1.ext", updateJob.SrcUrl)
	assert.EqualValues(t, "http://server/path2/file2.ext", updateJob.DstUrl)
	assert.EqualValues(t, "CreateAndRename", updateJob.Type)
	assert.EqualValues(t, "Created", updateJob.Status)
	assert.EqualValues(t, "abcdefg", updateJob.FileC4Id)
}

func TestUpdateJobSaveError(t *testing.T) {
	getJobFunction = func(jobId string) (*domain.Job, rest_errors.RestErr) {
		return &domain.Job{
			Id:         jobId,
			Name:       "Job 1",
			CreatedAt:  "2021-10-15T15:00:00Z",
			CreatedBy:  "user A",
			ModifiedAt: "",
			ModifiedBy: "",
			SrcUrl:     "http://server/path1/file1.ext",
			DstUrl:     "http://server/path2/file2.ext",
			Type:       "CreateAndRename",
			Status:     "Created",
			FileC4Id:   "abcdefg",
		}, nil
	}
	inputJob := domain.Job{
		Name:   "",
		DstUrl: "",
	}
	saveJobFunction = func(newJob domain.Job, overwrite bool) (*domain.Job, rest_errors.RestErr) {
		return nil, rest_errors.NewNotFoundError("could not save job")
	}
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	updateJob, err := JobService.Update(id, inputJob, true)
	assert.Nil(t, updateJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "could not save job", err.Message())
}

func TestGetNextNoJob(t *testing.T) {
	getNextJobFunction = func() (*domain.Job, rest_errors.RestErr) {
		return nil, rest_errors.NewNotFoundError("no jobs in list")
	}
	nextJob, err := JobService.GetNext()
	assert.Nil(t, nextJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "no jobs in list", err.Message())
}

func TestGetNextNoErro(t *testing.T) {
	getNextJobFunction = func() (*domain.Job, rest_errors.RestErr) {
		return &domain.Job{
			Id:         "1zXgBZNnBG1msmF1ARQK9ZphbbO",
			Name:       "Job 1",
			CreatedAt:  "2021-10-15T15:00:00Z",
			CreatedBy:  "user A",
			ModifiedAt: "",
			ModifiedBy: "",
			SrcUrl:     "http://server/path1/file1.ext",
			DstUrl:     "http://server/path2/file2.ext",
			Type:       "CreateAndRename",
			Status:     "Created",
			FileC4Id:   "abcdefg",
		}, nil
	}
	nextJob, err := JobService.GetNext()
	assert.NotNil(t, nextJob)
	assert.Nil(t, err)
	assert.EqualValues(t, "1zXgBZNnBG1msmF1ARQK9ZphbbO", nextJob.Id)
	assert.EqualValues(t, "Job 1", nextJob.Name)
}
