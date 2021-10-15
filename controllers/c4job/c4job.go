package c4job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/johannes-kuhfuss/c4/domain/c4job"
	services "github.com/johannes-kuhfuss/c4/services/C4JobService"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

func CreateC4Job(c *gin.Context) {
	var c4job domain.C4job
	if err := c.ShouldBindJSON(&c4job); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	result, err := services.C4JobService.CreateC4Job(c4job)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetC4Job(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented")
}
