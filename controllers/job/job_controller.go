package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/johannes-kuhfuss/c4/domain/job"
	services "github.com/johannes-kuhfuss/c4/services/JobService"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
	"github.com/segmentio/ksuid"
)

func getJobId(jobIdParam string) (string, rest_errors.RestErr) {
	jobId, err := ksuid.Parse(jobIdParam)
	if err != nil {
		return "", rest_errors.NewBadRequestError("user id should be a ksuid")
	}
	return jobId.String(), nil
}

func CreateJob(c *gin.Context) {
	var newJob domain.Job
	if err := c.ShouldBindJSON(&newJob); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	result, err := services.JobService.CreateJob(newJob)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetJob(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	job, err := services.JobService.GetJob(jobId)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, job)
}
