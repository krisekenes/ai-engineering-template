package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"example.com/ai-engineering-template/backend/internal/services"
)

type metadata struct {
	Timestamp string `json:"timestamp"`
}
type success[T any] struct {
	Data T        `json:"data"`
	Meta metadata `json:"meta"`
}
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type failure struct {
	Error apiError `json:"error"`
}

func New(now func() time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" && r.URL.Path != "/api/v1/estimate" {
			write(w, http.StatusNotFound, failure{apiError{"NOT_FOUND", "route not found"}})
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			write(w, http.StatusMethodNotAllowed, failure{apiError{"METHOD_NOT_ALLOWED", "use GET"}})
			return
		}
		meta := metadata{now().UTC().Format(time.RFC3339)}
		if r.URL.Path == "/healthz" {
			write(w, http.StatusOK, success[map[string]string]{map[string]string{"status": "ok"}, meta})
			return
		}
		values, present := r.URL.Query()["points"]
		if !present || len(values) != 1 {
			write(w, http.StatusBadRequest, failure{apiError{"VALIDATION_ERROR", services.ErrInvalidPoints.Error()}})
			return
		}
		points, err := strconv.Atoi(values[0])
		if err != nil {
			write(w, http.StatusBadRequest, failure{apiError{"VALIDATION_ERROR", services.ErrInvalidPoints.Error()}})
			return
		}
		estimate, err := services.Estimate(points)
		if err != nil {
			write(w, http.StatusBadRequest, failure{apiError{"VALIDATION_ERROR", err.Error()}})
			return
		}
		write(w, http.StatusOK, success[services.Estimation]{estimate, meta})
	})
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("write response failed", "error", err)
	}
}
