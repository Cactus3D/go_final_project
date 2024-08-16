package updatetask

import (
	"encoding/json"
	"net/http"

	"cactus3d/go_final_project/internal/models"
)

type TaskProvider interface {
	Update(*models.Task) error
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TaskProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req models.Task
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Некорректный формат запроса"})
			return
		}

		if ok, err := req.Validate(); !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		err = provider.Update(&req)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}
