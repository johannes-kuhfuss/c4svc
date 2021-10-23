package domain

import (
	"net/http"
	"testing"

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
)

func TestGetNotFound(t *testing.T) {
	user, err := JobDao.Get("X")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "job with Id X does not exist", err.Message())
}

func TestGetNoError(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	testJob, err := JobDao.Get("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, "1zXgBZNnBG1msmF1ARQK9ZphbbO", testJob.Id)
}

func TestSaveJobExistsNoOverwrite(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	newJob := Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbbO",
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
	assert.EqualValues(t, "job with Id 1zXgBZNnBG1msmF1ARQK9ZphbbO already exists", err.Message())
}

func TestSaveJobExistsOverwrite(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	newJob := Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbbO",
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
	testJob, err := JobDao.Save(newJob, true)
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, "1zXgBZNnBG1msmF1ARQK9ZphbbO", testJob.Id)
	assert.EqualValues(t, "Job 2", testJob.Name)
}

func TestDeleteJobNotFound(t *testing.T) {
	err := JobDao.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "job with Id 1zXgBZNnBG1msmF1ARQK9ZphbbO does not exist", err.Message())
}

func TestDeleteJobNoError(t *testing.T) {
	addJob(job1)
	defer removeJob(job1)
	err := JobDao.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
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
	assert.EqualValues(t, "1zXgBZNnBG1msmF1ARQK9ZphbdO", nextJob.Id)
	assert.EqualValues(t, "Job 3", nextJob.Name)
}
