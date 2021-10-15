package domain

import (
	"encoding/json"
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
	assert.EqualValues(t, JobStatusQueued, "Queued")
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
