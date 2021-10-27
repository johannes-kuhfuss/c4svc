package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/c4svc/utils/logger"
)

const (
	pong = "pong"
)

func Pong(c *gin.Context) {
	logger.Debug("Processing ping get request")
	c.String(http.StatusOK, pong)
	logger.Debug("Done processing ping get request")
}
