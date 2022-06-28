package server

import (
	"net/http"
	"os"

	"github.com/dawidl022/mooc-fi-kubernetes/todo/server/handlers"
	"gorm.io/gorm"
)

func routes(db *gorm.DB) {
	http.Handle("/api/hash", handlers.NewHashHandler())
	http.Handle("/api/daily-image", handlers.NewDailyImageHandler())
	http.Handle("/api/todos", handlers.NewTodoHandler(db))

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(0)
	})
}
