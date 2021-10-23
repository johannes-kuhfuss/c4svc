package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/johannes-kuhfuss/c4/domain/job"
	services "github.com/johannes-kuhfuss/c4/services/jobservice"
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

func Create(c *gin.Context) {
	var newJob domain.Job
	if err := c.ShouldBindJSON(&newJob); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	result, err := services.JobService.Create(newJob)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func Get(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	job, err := services.JobService.Get(jobId)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, job)
}

func Delete(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	err = services.JobService.Delete(jobId)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.String(http.StatusNoContent, "")
}

func validateUpdate(c *gin.Context) (id string, job domain.Job, err rest_errors.RestErr) {
	var inputJob domain.Job
	if err := c.ShouldBindJSON(&inputJob); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		return "", inputJob, apiErr
	}
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		return "", inputJob, err
	}
	return jobId, inputJob, nil
}

func Update(c *gin.Context) {
	partial := false
	jobId, inputJob, err := validateUpdate(c)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	result, err := services.JobService.Update(jobId, inputJob, partial)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func UpdatePart(c *gin.Context) {
	partial := true
	jobId, inputJob, err := validateUpdate(c)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	result, err := services.JobService.Update(jobId, inputJob, partial)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}
