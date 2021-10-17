package domain

import (
	"strings"

	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

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

type Job struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  string    `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ModifiedAt string    `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
	SrcUrl     string    `json:"src_url"`
	DstUrl     string    `json:"dst_url"`
	Type       JobType   `json:"type"`
	Status     JobStatus `json:"status"`
	FileC4Id   string    `json:"file_c4_id"`
}

func (j *Job) Validate() rest_errors.RestErr {
	if (j.Type != JobTypeCreate) && (j.Type != JobTypeCreateAndRename) {
		return rest_errors.NewBadRequestError("invalid job type")
	}
	if strings.TrimSpace(j.SrcUrl) == "" {
		return rest_errors.NewBadRequestError("invalid source Url")
	}
	if (j.Type == JobTypeCreateAndRename) && (strings.TrimSpace(j.DstUrl) == "") {
		return rest_errors.NewBadRequestError("invalid destination Url")
	}
	return nil
}