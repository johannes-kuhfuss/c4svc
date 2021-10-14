package c4job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//_ "github.com/johannes-kuhfuss/services_utils/rest_errors"
)

func CreateC4Job(c *gin.Context) {
	//restErr := rest_errors.NewBadRequestError("asdf")
	/*
		var c4job c4job.C4job
		if err := c.ShouldBindJSON(&c4job); err != nil {
			apiErr := restError.NewBadRequestError("invalid json body")
			c.JSON(apiErr.Status(), apiErr)
			return
		}

		clientId := c.GetHeader("X-Client-Id")

		result, err := services.RepositoryService.CreateRepo(clientId, request)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}
		c.JSON(http.StatusCreated, result)
	*/
}

func GetC4Job(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented")
}
