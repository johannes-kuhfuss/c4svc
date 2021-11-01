package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/johannes-kuhfuss/c4svc/domain/job"
	"github.com/johannes-kuhfuss/c4svc/services"
	"github.com/johannes-kuhfuss/c4svc/utils/logger"
	rest_errors "github.com/johannes-kuhfuss/c4svc/utils/rest_errors_utils"
	"github.com/segmentio/ksuid"
)

var (
	JobController jobControllerInterface = &jobController{}
)

type jobControllerInterface interface {
	Create(*gin.Context)
	Get(*gin.Context)
	Delete(*gin.Context)
	Update(*gin.Context)
	UpdatePart(*gin.Context)
	GetAll(*gin.Context)
}

type jobController struct {
}

func getJobId(jobIdParam string) (string, rest_errors.RestErr) {
	jobId, err := ksuid.Parse(jobIdParam)
	if err != nil {
		logger.Error("User Id should be a ksuid", err)
		return "", rest_errors.NewBadRequestError("user id should be a ksuid")
	}
	return jobId.String(), nil
}

func (jc jobController) Create(c *gin.Context) {
	logger.Debug("Processing job create request")
	var newJob domain.Job
	if err := c.ShouldBindJSON(&newJob); err != nil {
		logger.Error("invalid JSON body in create request", err)
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	result, err := services.JobService.Create(newJob)
	if err != nil {
		logger.Error("Service error while creating job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
	logger.Debug("Done processing job create request")
}

func (jc jobController) Get(c *gin.Context) {
	logger.Debug("Processing job get request")
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	job, err := services.JobService.Get(jobId)
	if err != nil {
		logger.Error("Service error while getting job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, job)
	logger.Debug("Done processing job get request")
}

func (jc jobController) Delete(c *gin.Context) {
	logger.Debug("Processing job delete request")
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	err = services.JobService.Delete(jobId)
	if err != nil {
		logger.Error("Service error while deleting job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.String(http.StatusNoContent, "")
	logger.Debug("Done processing job delete request")
}

func validateUpdate(c *gin.Context) (id string, job domain.Job, err rest_errors.RestErr) {
	logger.Debug("Validating update")
	var inputJob domain.Job
	if err := c.ShouldBindJSON(&inputJob); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		return "", inputJob, apiErr
	}
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		return "", inputJob, err
	}
	logger.Debug("Done validating update")
	return jobId, inputJob, nil
}

func (jc jobController) Update(c *gin.Context) {
	logger.Debug("Processing job full update request")
	partial := false
	jobId, inputJob, err := validateUpdate(c)
	if err != nil {
		logger.Error("Error while validating full job update", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	result, err := services.JobService.Update(jobId, inputJob, partial)
	if err != nil {
		logger.Error("Service error while updating full job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, result)
	logger.Debug("Done processing job full update request")
}

func (jc jobController) UpdatePart(c *gin.Context) {
	logger.Debug("Processing job partial update request")
	partial := true
	jobId, inputJob, err := validateUpdate(c)
	if err != nil {
		logger.Error("Error while validating partial job update", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	result, err := services.JobService.Update(jobId, inputJob, partial)
	if err != nil {
		logger.Error("Service error while updating partial job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, result)
	logger.Debug("Done processing job partial update request")
}

func (jc jobController) GetAll(c *gin.Context) {
	logger.Debug("Processing job get all request")
	jobs, err := services.JobService.GetAll()
	if err != nil {
		logger.Error("Service error while getting all jobs", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, jobs)
	logger.Debug("Done processing job get request")
}
