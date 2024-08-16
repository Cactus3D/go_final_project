package gettask

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cactus3d/go_final_project/internal/models"
)

type TaskProvider interface {
	Get(id string) (*models.Task, error)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TaskProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсуетствует id"})
			return
		}
		if id, err := strconv.Atoi(id); err != nil || id < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "id должен быть положительным числом"})
			return
		}

		res, err := provider.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсутвует task"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*res)
	}
}
