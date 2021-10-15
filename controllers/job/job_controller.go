package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/johannes-kuhfuss/c4/domain/job"
	services "github.com/johannes-kuhfuss/c4/services/JobService"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

func CreateJob(c *gin.Context) {
	var newJob domain.Job
	if err := c.ShouldBindJSON(&newJob); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	//fmt.Printf("Job: %#v", newJob)

	result, err := services.JobService.CreateJob(newJob)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetJob(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented")
}
