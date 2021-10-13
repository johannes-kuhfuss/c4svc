package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	pong = "pong"
)

func Pong(c *gin.Context) {
	c.String(http.StatusOK, pong)
}
