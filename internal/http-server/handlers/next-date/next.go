package nextdate

import (
	"net/http"
	"time"

	"cactus3d/go_final_project/internal/nextdate"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		now := r.URL.Query().Get("now")
		date := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")

		n, err := time.Parse(nextdate.DateFormat, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = time.Parse(nextdate.DateFormat, date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := nextdate.NextDate(n, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(res))
	}
}
