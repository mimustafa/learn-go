package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func countRequests(next http.Handler) http.Handler {
	count := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.Header().Set("X-Request-Count", strconv.Itoa(count))
		next.ServeHTTP(w, r)
	})
}

func main() {
	observations := []string{"heart-rate: 72 bpm", "sleep-hours: 8 h"}
	mux := http.NewServeMux()
	mux.HandleFunc("/observations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
			var request struct {
				Code string `json:"code"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&request); err != nil {
				writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
				return
			}
			if strings.TrimSpace(request.Code) == "" {
				writeError(w, http.StatusBadRequest, "invalid_code", "code is required")
				return
			}
			if err := decoder.Decode(new(any)); err != io.EOF {
				writeError(w, http.StatusBadRequest, "invalid_json", "request body must contain one JSON value")
				return
			}
			observations = append(observations, request.Code)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": request.Code})
			return
		case http.MethodGet:
			limit, offset := 2, 0
			var err error
			if raw := r.URL.Query().Get("limit"); raw != "" {
				limit, err = strconv.Atoi(raw)
				if err != nil || limit < 1 {
					writeError(w, http.StatusBadRequest, "invalid_limit", "limit must be a positive integer")
					return
				}
			}
			if limit > 50 {
				limit = 50
			}
			if raw := r.URL.Query().Get("offset"); raw != "" {
				offset, err = strconv.Atoi(raw)
				if err != nil || offset < 0 {
					writeError(w, http.StatusBadRequest, "invalid_offset", "offset must be a nonnegative integer")
					return
				}
			}
			if offset > len(observations) {
				offset = len(observations)
			}
			end := offset + limit
			if end > len(observations) {
				end = len(observations)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"observations": observations[offset:end], "limit": limit, "offset": offset,
			})
			return
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	handler := countRequests(mux)
	invalid := httptest.NewRecorder()
	handler.ServeHTTP(invalid, httptest.NewRequest(http.MethodPost, "/observations", strings.NewReader(`{"code":" "}`)))
	fmt.Printf("invalid observation status %d, body %s", invalid.Code, invalid.Body.String())
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/observations?limit=100&offset=1", nil))
	fmt.Printf("status %d, count %s, body %s", recorder.Code, recorder.Header().Get("X-Request-Count"), recorder.Body.String())
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError{Code: code, Message: message})
}
