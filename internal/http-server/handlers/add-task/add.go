package addtask

import (
	"encoding/json"
	"net/http"
	"time"

	"cactus3d/go_final_project/internal/nextdate"
)

type TasksProvider interface {
	Add(date, title, comment, repeat string) (int, error)
}

type Request struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	Id int `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(provider TasksProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Некорректный формат запроса"})
			return
		}

		if req.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Не указан заголовок задачи"})
			return
		}

		if req.Date == "" {
			req.Date = time.Now().Format(nextdate.DateFormat)
		} else {
			_, err = time.Parse(nextdate.DateFormat, req.Date)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат времени"})
				return
			}
		}

		if req.Repeat != "" {
			_, err = nextdate.NextDate(time.Now(), req.Date, req.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат повторений"})
				return
			}
		}

		id, err := provider.Add(req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Id: id})
	}
}
