package gettasks

import (
	"encoding/json"
	"net/http"

	"cactus3d/go_final_project/internal/models"
)

type TaskProvider interface {
	GetAll(search string) ([]models.Task, error)
}

type Response struct {
	Tasks []models.Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TaskProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		var err error

		search := r.URL.Query().Get("search")

		tasks, err := provider.GetAll(search)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Tasks: tasks})
	}
}
