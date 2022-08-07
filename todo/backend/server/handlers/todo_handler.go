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
		t.get(w)
	case http.MethodPost:
		t.post(w, r)
	case http.MethodPut:
		t.put(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func (t *todoHandler) get(w http.ResponseWriter) {
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
}

func (t *todoHandler) post(w http.ResponseWriter, r *http.Request) {
	todo, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to get todo from request", http.StatusInternalServerError)
		return
	}
	if len(todo) > 140 {
		http.Error(w, "Todo cannot be longer than 140 characters", http.StatusBadRequest)
		return
	}
	todoModel := models.Todo{Content: string(todo)}
	err = t.db.Create(&todoModel).Error
	if err != nil {
		http.Error(w, "Failed to save todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t *todoHandler) put(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to get todo from request", http.StatusInternalServerError)
		return
	}

	var todo models.Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		http.Error(w, "Failed to parse todo", http.StatusBadRequest)
		return
	}
	if err := t.db.Save(&todo).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AddWikiPage(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get("https://en.wikipedia.org/wiki/Special:Random")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := resp.Header.Get("Location")
		todo := models.Todo{Content: url}

		err = db.Create(&todo).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
