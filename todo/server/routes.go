package server

import (
	"net/http"

	"github.com/dawidl022/mooc-fi-kubernetes/todo/server/handlers"
)

func routes() {
	http.Handle("/hash", handlers.NewHashHandler())
	http.Handle("/daily-image", handlers.NewDailyImageHandler())
}
