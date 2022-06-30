package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dawidl022/mooc-fi-kubernetes/todo/models"
	"gorm.io/gorm"
)

type todoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *todoHandler {
	return &todoHandler{
		db: db,
	}
}

func (t *todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		var todos []*models.Todo
		err := t.db.Find(&todos).Error
		if err != nil {
			http.Error(w, "failed to fetch todos", http.StatusInternalServerError)
			return
		}

		todosJson, err := json.Marshal(todos)
		if err != nil {
			http.Error(w, "failed to marshal todos", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(todosJson)
	case http.MethodPost:
		todo, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to get todo from request", http.StatusInternalServerError)
			return
		}
		todoModel := models.Todo{Content: string(todo)}
		err = t.db.Create(&todoModel).Error
		if err != nil {
			http.Error(w, "Failed to save todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func AddWikiPage(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://en.wikipedia.org/wiki/Special:Random")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := resp.Header.Get("location")
		todo := models.Todo{Content: url}

		err = db.Create(&todo).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
