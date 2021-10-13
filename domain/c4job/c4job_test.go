package c4job

import (
	"encoding/json"
	"testing"

	"github.com/johannes-kuhfuss/c4/utils/date_utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateC4JobAsJson(t *testing.T) {
	job1 := C4job{
		Id:         1,
		Name:       "new C4 Job",
		CreatedAt:  date_utils.GetNowUtcString(),
		CreatedBy:  "user1",
		ModifiedAt: date_utils.GetNowUtcString(),
		ModifiedBy: "user2",
		SrcUrl:     "https://server/path1/file1.ext",
		DstUrl:     "https://server/path2/file2.ext",
		//Type:       JobType.Create,
		//Status:     JobStatus.Created,
	}
	bytes, err := json.Marshal(job1)
	assert.NotNil(t, bytes)
	assert.Nil(t, err)

	var job2 C4job
	err = json.Unmarshal(bytes, &job2)
	assert.Nil(t, err)
	assert.NotNil(t, job2)
	assert.EqualValues(t, job2, job1)
}
