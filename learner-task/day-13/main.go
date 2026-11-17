package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errorJSON(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError{Code: code, Message: message})
}

func withRequestID(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a fresh ID rather than trusting an unchecked client header.
		var bytes [8]byte
		if _, err := rand.Read(bytes[:]); err != nil {
			errorJSON(w, http.StatusInternalServerError, "request_id_failed", "request could not be served")
			return
		}
		id := hex.EncodeToString(bytes[:])
		w.Header().Set("X-Request-ID", id)
		logger.Info("CareLog request", "request_id", id, "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func paged(items []string, rawLimit, rawOffset string) ([]string, int, int, error) {
	limit, offset := 20, 0
	var err error
	if rawLimit != "" {
		limit, err = strconv.Atoi(rawLimit)
		if err != nil || limit < 1 {
			return nil, 0, 0, fmt.Errorf("limit must be positive")
		}
	}
	if limit > 50 {
		// Cap client-selected work to keep responses and queries bounded.
		limit = 50
	}
	if rawOffset != "" {
		offset, err = strconv.Atoi(rawOffset)
		if err != nil || offset < 0 {
			return nil, 0, 0, fmt.Errorf("offset must be nonnegative")
		}
	}
	if offset >= len(items) {
		return []string{}, limit, offset, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], limit, offset, nil
}

func handler() http.Handler {
	observations := []string{"heart-rate", "sleep-hours", "steps"}
	mux := http.NewServeMux()
	mux.HandleFunc("/observations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
			var body struct {
				Code string `json:"code"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&body); err != nil {
				errorJSON(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
				return
			}
			if strings.TrimSpace(body.Code) == "" {
				errorJSON(w, http.StatusUnprocessableEntity, "code_required", "observation code is required")
				return
			}
			if len(body.Code) > 64 {
				errorJSON(w, http.StatusUnprocessableEntity, "code_too_long", "observation code exceeds 64 bytes")
				return
			}
			observations = append(observations, body.Code)
			w.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			items, limit, offset, err := paged(observations, r.URL.Query().Get("limit"), r.URL.Query().Get("offset"))
			if err != nil {
				errorJSON(w, http.StatusBadRequest, "invalid_page", err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"items": items, "limit": limit, "offset": offset})
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return withRequestID(logger, mux)
}

func main() {
	h := handler()
	for _, request := range []*http.Request{
		httptest.NewRequest(http.MethodPost, "/observations", strings.NewReader(`{"code":" "}`)),
		httptest.NewRequest(http.MethodGet, "/observations?limit=500&offset=2", nil),
	} {
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		fmt.Printf("status=%d request_id=%s body=%s", recorder.Code, recorder.Header().Get("X-Request-ID"), recorder.Body.String())
	}
}
