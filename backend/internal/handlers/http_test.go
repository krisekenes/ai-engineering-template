package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/ai-engineering-template/backend/internal/services"
)

func TestHTTPContract(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.FixedZone("offset", 3600))
	handler := New(func() time.Time { return now })
	for _, tc := range []struct {
		name, method, path string
		status             int
		code               string
	}{
		{"success", "GET", "/api/v1/estimate?points=3", 200, ""},
		{"minimum", "GET", "/api/v1/estimate?points=1", 200, ""},
		{"maximum", "GET", "/api/v1/estimate?points=100", 200, ""},
		{"missing", "GET", "/api/v1/estimate", 400, "VALIDATION_ERROR"},
		{"empty", "GET", "/api/v1/estimate?points=", 400, "VALIDATION_ERROR"},
		{"negative", "GET", "/api/v1/estimate?points=-1", 400, "VALIDATION_ERROR"},
		{"zero", "GET", "/api/v1/estimate?points=0", 400, "VALIDATION_ERROR"},
		{"overflow range", "GET", "/api/v1/estimate?points=101", 400, "VALIDATION_ERROR"},
		{"overflow integer", "GET", "/api/v1/estimate?points=99999999999999999999999999", 400, "VALIDATION_ERROR"},
		{"decimal", "GET", "/api/v1/estimate?points=1.5", 400, "VALIDATION_ERROR"},
		{"text", "GET", "/api/v1/estimate?points=hello", 400, "VALIDATION_ERROR"},
		{"duplicate", "GET", "/api/v1/estimate?points=1&points=2", 400, "VALIDATION_ERROR"},
		{"method", "POST", "/api/v1/estimate?points=3", 405, "METHOD_NOT_ALLOWED"},
		{"health method", "POST", "/healthz", 405, "METHOD_NOT_ALLOWED"},
		{"unknown", "GET", "/missing", 404, "NOT_FOUND"},
		{"health", "GET", "/healthz", 200, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.status {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
			}
			if rec.Header().Get("Content-Type") != "application/json" {
				t.Fatal("missing JSON content type")
			}
			if tc.status == http.StatusMethodNotAllowed && rec.Header().Get("Allow") != "GET" {
				t.Fatal("missing Allow header")
			}
			if tc.code != "" {
				var got failure
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.Error.Code != tc.code || got.Error.Message == "" {
					t.Fatalf("error = %+v", got)
				}
			} else if tc.name == "health" {
				var got success[map[string]string]
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.Data["status"] != "ok" {
					t.Fatalf("health = %+v", got)
				}
			} else {
				var got success[services.Estimation]
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.Meta.Timestamp != "2026-01-01T11:00:00Z" {
					t.Fatalf("timestamp = %s", got.Meta.Timestamp)
				}
				if got.Data.Hours != got.Data.Points*4 || got.Data.Points < 1 {
					t.Fatalf("estimate = %+v", got.Data)
				}
				if tc.name == "success" && (got.Data.Points != 3 || got.Data.Hours != 12) {
					t.Fatalf("estimate = %+v", got.Data)
				}
			}
		})
	}
}
