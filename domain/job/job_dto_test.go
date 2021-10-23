package domain

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/johannes-kuhfuss/c4/utils/date_utils"
	"github.com/stretchr/testify/assert"
)

func TestConstJobType(t *testing.T) {
	assert.EqualValues(t, JobTypeCreate, "Create")
	assert.EqualValues(t, JobTypeCreateAndRename, "CreateAndRename")
}

func TestConstJobStatus(t *testing.T) {
	assert.EqualValues(t, JobStatusCreated, "Created")
	assert.EqualValues(t, JobStatusRunning, "Running")
	assert.EqualValues(t, JobStatusFinished, "Finished")
	assert.EqualValues(t, JobStatusFailed, "Failed")
}

func TestCreateC4JobAsJson(t *testing.T) {
	job1 := Job{
		Id:         "1zXg7ubtb02J3t0muj6jXqzzM72",
		Name:       "new C4 Job",
		CreatedAt:  date_utils.GetNowUtcString(),
		CreatedBy:  "user1",
		ModifiedAt: date_utils.GetNowUtcString(),
		ModifiedBy: "user2",
		SrcUrl:     "https://server/path1/file1.ext",
		DstUrl:     "https://server/path2/file2.ext",
		Type:       JobTypeCreate,
		Status:     JobStatusCreated,
		FileC4Id:   "abcdefg",
	}
	bytes, err := json.Marshal(job1)
	assert.NotNil(t, bytes)
	assert.Nil(t, err)

	var job2 Job
	err = json.Unmarshal(bytes, &job2)
	assert.Nil(t, err)
	assert.NotNil(t, job2)
	assert.EqualValues(t, job2, job1)
}

func TestValidateWrongJobType(t *testing.T) {
	job1 := Job{
		Type: "not correct",
	}
	err := job1.Validate()
	assert.NotNil(t, err)
	assert.EqualValues(t, err.StatusCode(), http.StatusBadRequest)
	assert.EqualValues(t, err.Message(), "invalid job type")
}

func TestValidateNoSource(t *testing.T) {
	job1 := Job{
		Type: "Create",
	}
	err := job1.Validate()
	assert.NotNil(t, err)
	assert.EqualValues(t, err.StatusCode(), http.StatusBadRequest)
	assert.EqualValues(t, err.Message(), "invalid source Url")
}

func TestValidateRenameNoDestination(t *testing.T) {
	job1 := Job{
		Type:   "CreateAndRename",
		SrcUrl: "https://server/path1/file1.ext",
	}
	err := job1.Validate()
	assert.NotNil(t, err)
	assert.EqualValues(t, err.StatusCode(), http.StatusBadRequest)
	assert.EqualValues(t, err.Message(), "invalid destination Url")
}

func TestValidateNoError(t *testing.T) {
	job1 := Job{
		Type:   "CreateAndRename",
		SrcUrl: "https://server/path1/file1.ext",
		DstUrl: "https://server/path2/file2.ext",
	}
	err := job1.Validate()
	assert.Nil(t, err)
}
