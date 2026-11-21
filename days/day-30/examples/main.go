package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
)

type Observation struct {
	ID        string  `json:"id"`
	PatientID string  `json:"patient_id"`
	Code      string  `json:"code"`
	Value     float64 `json:"value"`
}

type createRequest struct {
	PatientID string   `json:"patient_id"`
	Code      string   `json:"code"`
	Value     *float64 `json:"value"`
}

type Store struct {
	mu     sync.Mutex
	byKey  map[string]Observation
	nextID uint64
}

func (s *Store) Create(key string, input createRequest) (Observation, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if input.Value == nil {
		return Observation{}, false, errors.New("value is required")
	}
	if old, ok := s.byKey[key]; ok {
		if old.PatientID != input.PatientID || old.Code != input.Code || old.Value != *input.Value {
			return Observation{}, true, errors.New("idempotency key reused with different input")
		}
		return old, true, nil
	}
	s.nextID++
	item := Observation{
		ID: fmt.Sprintf("obs-%d", s.nextID), PatientID: input.PatientID,
		Code: input.Code, Value: *input.Value,
	}
	s.byKey[key] = item
	return item, false, nil
}

type requestIDKey struct{}

func newHandler(store *Store) http.Handler {
	var requestNumber atomic.Uint64
	mux := http.NewServeMux()
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/observations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		requestID := fmt.Sprintf("req-%d", requestNumber.Add(1))
		w.Header().Set("X-Request-ID", requestID)
		r = r.WithContext(context.WithValue(r.Context(), requestIDKey{}, requestID))

		key := r.Header.Get("Idempotency-Key")
		if key == "" || len(key) > 80 {
			http.Error(w, "valid idempotency key required", http.StatusBadRequest)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var input createRequest
		if err := decoder.Decode(&input); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			http.Error(w, "one JSON object required", http.StatusBadRequest)
			return
		}
		if !validIdentifier(input.PatientID) || !validCode(input.Code) {
			http.Error(w, "invalid observation fields", http.StatusBadRequest)
			return
		}
		if input.Value == nil || math.IsNaN(*input.Value) || math.IsInf(*input.Value, 0) ||
			*input.Value < 0 || *input.Value > 1000 {
			http.Error(w, "value out of range", http.StatusBadRequest)
			return
		}
		item, duplicate, err := store.Create(key, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		correlationID, _ := r.Context().Value(requestIDKey{}).(string)
		slog.InfoContext(r.Context(), "observation accepted",
			"request_id", correlationID, "duplicate", duplicate)
		w.Header().Set("Content-Type", "application/json")
		if duplicate {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		_ = json.NewEncoder(w).Encode(item)
	})
	return mux
}

func validIdentifier(value string) bool {
	if len(value) < 1 || len(value) > 40 {
		return false
	}
	for _, r := range value {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' ||
			r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

func validCode(value string) bool {
	if len(value) < 1 || len(value) > 12 {
		return false
	}
	for _, r := range value {
		if !(r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func main() {
	store := &Store{byKey: make(map[string]Observation)}
	server := httptest.NewServer(newHandler(store))
	defer server.Close()

	body := `{"patient_id":"patient-1","code":"TEMP","value":37.2}`
	for attempt := 0; attempt < 2; attempt++ {
		request, _ := http.NewRequest(http.MethodPost, server.URL+"/observations", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "submit-001")
		response, err := server.Client().Do(request)
		if err != nil {
			panic(err)
		}
		fmt.Printf("attempt=%d status=%d request_id=%s\n",
			attempt+1, response.StatusCode, response.Header.Get("X-Request-ID"))
		response.Body.Close()
	}
}
