package c4job

type JobType string

const (
	JobTypeCreate          = "Create"
	JobTypeCreateAndRename = "CreateAndRename"
)

type JobStatus string

const (
	JobStatusCreated  = "Created"
	JobStatusQueued   = "Queued"
	JobStatusRunning  = "Running"
	JobStatusFinished = "Finished"
	JobStatusFailed   = "Failed"
)

type C4job struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  string    `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ModifiedAt string    `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
	SrcUrl     string    `json:"src_url"`
	DstUrl     string    `json:"dst_url"`
	Type       JobType   `json:"job_type"`
	Status     JobStatus `json:"status"`
}
