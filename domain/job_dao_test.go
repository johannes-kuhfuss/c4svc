package domain

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/johannes-kuhfuss/c4svc/config"
	"github.com/johannes-kuhfuss/c4svc/utils/date"
	"github.com/stretchr/testify/assert"
)

var (
	job1 Job = Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbbO",
		Name:       "Job 1",
		CreatedAt:  "2021-10-15T15:00:00Z",
		CreatedBy:  "user 1",
		ModifiedAt: "",
		ModifiedBy: "",
		SrcUrl:     "http://server1/path1/file1.ext",
		DstUrl:     "",
		Type:       "Create",
		Status:     "Running",
		FileC4Id:   "abcdefg",
	}
	job2 Job = Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbcO",
		Name:       "Job 2",
		CreatedAt:  "2021-10-15T15:00:00Z",
		CreatedBy:  "user 2",
		ModifiedAt: "",
		ModifiedBy: "",
		SrcUrl:     "http://server2/path2/file2.ext",
		DstUrl:     "",
		Type:       "Create",
		Status:     "Created",
		FileC4Id:   "abcdefg",
	}
	job3 Job = Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbdO",
		Name:       "Job 3",
		CreatedAt:  "2021-10-14T14:00:00Z",
		CreatedBy:  "user 3",
		ModifiedAt: "",
		ModifiedBy: "",
		SrcUrl:     "http://server3/path3/file3.ext",
		DstUrl:     "",
		Type:       "Create",
		Status:     "Created",
		FileC4Id:   "abcdefg",
	}
	job4 Job = Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbdO",
		Name:       "Job 3",
		CreatedAt:  "2021-10-14T14:00:00Z",
		CreatedBy:  "user 3",
		ModifiedAt: "",
		ModifiedBy: "",
		SrcUrl:     "http://server3/path3/file3.ext",
		DstUrl:     "",
		Type:       "Create",
		Status:     "Finished",
		FileC4Id:   "abcdefg",
	}
)

func TestGetNotFound(t *testing.T) {
	id := "X"
	user, err := JobDao.Get(id)
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("job with Id %v does not exist", id), err.Message())
}

func TestGetNoError(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	testJob, err := JobDao.Get(job1.Id)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, job1.Id, testJob.Id)
}

func TestSaveJobExistsNoOverwrite(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	newJob := Job{
		Id:         fmt.Sprintf("%v", id),
		Name:       "Job 2",
		CreatedAt:  "2021-10-15T15:00:00Z",
		CreatedBy:  "user A",
		ModifiedAt: "",
		ModifiedBy: "",
		SrcUrl:     "http://server/path1/file1.ext",
		DstUrl:     "",
		Type:       "Create",
		Status:     "Running",
		FileC4Id:   "abcdefg",
	}
	testJob, err := JobDao.Save(newJob, false)
	assert.Nil(t, testJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("job with Id %v already exists", id), err.Message())
}

func TestSaveJobExistsOverwrite(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	name := "Job 2"
	newJob := Job{
		Id:         fmt.Sprintf("%v", id),
		Name:       fmt.Sprintf("%v", name),
		CreatedAt:  "2021-10-15T15:00:00Z",
		CreatedBy:  "user A",
		ModifiedAt: "",
		ModifiedBy: "",
		SrcUrl:     "http://server/path1/file1.ext",
		DstUrl:     "",
		Type:       "Create",
		Status:     "Running",
		FileC4Id:   "abcdefg",
	}
	testJob, err := JobDao.Save(newJob, true)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, id, testJob.Id)
	assert.EqualValues(t, name, testJob.Name)
}

func TestDeleteJobNotFound(t *testing.T) {
	id := "1zXgBZNnBG1msmF1ARQK9ZphbbO"
	err := JobDao.Delete(id)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("job with Id %v does not exist", id), err.Message())
}

func TestDeleteJobNoError(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.Delete(job1.Id)
	assert.Nil(t, err)
}

func TestGetNextListEmpty(t *testing.T) {
	nextJob, err := JobDao.GetNext()
	assert.Nil(t, nextJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "no jobs in list", err.Message())
}

func TestGetNextNoError(t *testing.T) {
	addJob(job1)
	addJob(job2)
	addJob(job3)
	defer removeJob(job1)
	defer removeJob(job2)
	defer removeJob(job3)
	nextJob, err := JobDao.GetNext()
	assert.NotNil(t, nextJob)
	assert.Nil(t, err)
	assert.EqualValues(t, job3.Id, nextJob.Id)
	assert.EqualValues(t, job3.Name, nextJob.Name)
}

func TestChangeStatusNoJob(t *testing.T) {
	id := "1zXgBZNnBG1msmF1ARQK9ZphbdO"
	err := JobDao.ChangeStatus(id, "")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("job with Id %v does not exist", id), err.Message())
}

func TestChangeStatusInvalidStatus(t *testing.T) {
	addJob(job3)
	defer removeJob(job3)
	err := JobDao.ChangeStatus(job3.Id, "invalidstatus")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "invalid status value", err.Message())
}

func TestChangeStatusSameStatus(t *testing.T) {
	addJob(job3)
	defer removeJob(job3)
	err := JobDao.ChangeStatus(job3.Id, "created")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job3.Id)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatus("Created"), testJob.Status)
}

func TestChangeStatusNoErrorCreated(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.ChangeStatus(job1.Id, "created")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job1.Id)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatus("Created"), testJob.Status)
}

func TestChangeStatusNoErrorRunning(t *testing.T) {
	addJob(job3)
	defer removeJob(job3)
	err := JobDao.ChangeStatus(job3.Id, "running")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job3.Id)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatus("Running"), testJob.Status)
}

func TestChangeStatusNoErrorFailed(t *testing.T) {
	addJob(job3)
	defer removeJob(job3)
	err := JobDao.ChangeStatus(job3.Id, "failed")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job3.Id)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatus("Failed"), testJob.Status)
}

func TestChangeStatusNoErrorFinished(t *testing.T) {
	addJob(job3)
	defer removeJob(job3)
	err := JobDao.ChangeStatus(job3.Id, "finished")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job3.Id)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatus("Finished"), testJob.Status)
}

func TestCleanJobsNoJobs(t *testing.T) {
	numJobs, err := JobDao.CleanJobs(config.DeleteFinishedAge, config.DeleteFailedAge)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, numJobs)
}

func TestCleanJobsNoModDate(t *testing.T) {
	addJob(job4)
	numJobs, err := JobDao.CleanJobs(config.DeleteFinishedAge, config.DeleteFailedAge)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, numJobs)
}

func TestCleanJobsNoError(t *testing.T) {
	job4.ModifiedAt = date.GetNowUtc().Add(-config.DeleteFinishedAge).Format(date.ApiDateLayout)
	addJob(job4)
	numJobs, err := JobDao.CleanJobs(config.DeleteFinishedAge, config.DeleteFailedAge)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, numJobs)
}

func TestSetC4IdNoJobFound(t *testing.T) {
	id := "1zXgBZNnBG1msmF1ARQK9ZphbdO"
	err := JobDao.SetC4Id(id, "C4id")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("job with Id %v does not exist", id), err.Message())
}

func TestSetC4IdInvalidC4Id(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.SetC4Id(job1.Id, "")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "invalid C4 Id", err.Message())
}

func TestSetC4IdNoError(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.SetC4Id(job1.Id, "new C4 Id")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job1.Id)
	assert.Nil(t, err)
	assert.EqualValues(t, "new C4 Id", testJob.FileC4Id)
}

func TestSetDstUrlNoJobFound(t *testing.T) {
	id := "1zXgBZNnBG1msmF1ARQK9ZphbdO"
	err := JobDao.SetDstUrl(id, "new url")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("job with Id %v does not exist", id), err.Message())
}

func TestSetDstUrlInvalidDstUrl(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.SetDstUrl(job1.Id, "")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "invalid destination URL", err.Message())
}

func TestSetDstUrlNoError(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.SetDstUrl(job1.Id, "new destination URL")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job1.Id)
	assert.Nil(t, err)
	assert.EqualValues(t, "new destination URL", testJob.DstUrl)
}

func TestSetErrMsgNoJobFound(t *testing.T) {
	id := "1zXgBZNnBG1msmF1ARQK9ZphbdO"
	err := JobDao.SetErrMsg(id, "new error message")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("job with Id %v does not exist", id), err.Message())
}

func TestSetErrMsgNoError(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.SetErrMsg(job1.Id, "new error message")
	assert.Nil(t, err)
	testJob, err := JobDao.Get(job1.Id)
	assert.Nil(t, err)
	assert.EqualValues(t, "new error message", testJob.ErrorMsg)
}

func TestGetAllNoJobsError(t *testing.T) {
	jobs, err := JobDao.GetAll()
	assert.Nil(t, jobs)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "no jobs in list", err.Message())
}

func TestGetAllNoError(t *testing.T) {
	addJob(job1)
	addJob(job2)
	addJob(job3)
	defer removeJob(job1)
	defer removeJob(job2)
	defer removeJob(job3)
	jobs, err := JobDao.GetAll()
	assert.NotNil(t, jobs)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(*jobs))
}
