package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type todoHandler struct {
	todos []string
}

func NewTodoHandler() *todoHandler {
	return &todoHandler{
		todos: []string{},
	}
}

func (t *todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		todos, err := json.Marshal(t.todos)
		if err != nil {
			http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(todos)
	case http.MethodPost:
		todo, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to save todo", http.StatusInternalServerError)
			return
		}
		t.todos = append(t.todos, string(todo))
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
