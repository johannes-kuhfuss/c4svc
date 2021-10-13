package app

import (
	"github.com/johannes-kuhfuss/c4/controllers/ping"
)

func mapUrls() {
	router.GET("/ping", ping.Pong)
}
