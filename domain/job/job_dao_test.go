package domain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestJob() {
	newJob := Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbbO",
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
	}
	addJob(newJob)
}

func deleteTestJob() {
	deleteJob := Job{
		Id:         "1zXgBZNnBG1msmF1ARQK9ZphbbO",
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
	}
	removeJob(deleteJob)
}

func TestGetNotFound(t *testing.T) {
	user, err := JobDao.Get("X")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "job with Id X does not exist", err.Message())
}

func TestGetNoError(t *testing.T) {
	setupTestJob()
	defer deleteTestJob()
	testJob, err := JobDao.Get("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.NotNil(t, testJob)
	assert.Nil(t, err)
	assert.EqualValues(t, "1zXgBZNnBG1msmF1ARQK9ZphbbO", testJob.Id)
}

func TestSaveJobExistsNoOverwrite(t *testing.T) {
	setupTestJob()
	defer deleteTestJob()
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
	setupTestJob()
	defer deleteTestJob()
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
	setupTestJob()
	defer deleteTestJob()
	err := JobDao.Delete("1zXgBZNnBG1msmF1ARQK9ZphbbO")
	assert.Nil(t, err)
}
