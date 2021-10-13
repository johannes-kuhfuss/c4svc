package c4job

type JobType int64

const (
	Create          JobType = 1
	CreateAndRename JobType = 2
)

func (jt JobType) String() string {
	switch jt {
	case Create:
		return "Create"
	case CreateAndRename:
		return "CreateAndRename"
	}
	return "unknown"
}

type JobStatus int64

const (
	Created  JobStatus = 1
	Queued   JobStatus = 2
	Running  JobStatus = 3
	Finished JobStatus = 4
	Failed   JobStatus = 5
)

func (js JobStatus) String() string {
	switch js {
	case Created:
		return "Created"
	case Queued:
		return "Queued"
	case Running:
		return "Running"
	case Finished:
		return "Finished"
	case Failed:
		return "Failed"
	}
	return "unknown"
}

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
