package server

import (
	"net/http"
	"os"

	"github.com/dawidl022/mooc-fi-kubernetes/todo/server/handlers"
	"gorm.io/gorm"
)

func routes(router *http.ServeMux, db *gorm.DB) {
	router.Handle("/api/hash", handlers.NewHashHandler())
	router.Handle("/api/daily-image", handlers.NewDailyImageHandler())
	router.Handle("/api/todos", handlers.NewTodoHandler(db))

	router.HandleFunc("/add-wiki-page", handlers.AddWikiPage(db))
	router.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(0)
	})
	router.HandleFunc("/healthz", handlers.PingDB(db))
}
