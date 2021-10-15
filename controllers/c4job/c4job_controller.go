package c4job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/johannes-kuhfuss/c4/domain/c4job"
	services "github.com/johannes-kuhfuss/c4/services/C4JobService"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

func CreateC4Job(c *gin.Context) {
	var newJob domain.C4job
	if err := c.ShouldBindJSON(&newJob); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	//fmt.Printf("Job: %#v", newJob)

	result, err := services.C4JobService.CreateC4Job(newJob)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetC4Job(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented")
}
