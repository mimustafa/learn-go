package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type result struct {
	ID string `json:"id"`
}

func observationService() (*httptest.Server, *atomic.Int64) {
	var mu sync.Mutex
	seen := make(map[string]result)
	var effects atomic.Int64
	mux := http.NewServeMux()
	mux.HandleFunc("/observations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "idempotency key required", http.StatusBadRequest)
			return
		}
		mu.Lock()
		item, duplicate := seen[key]
		if !duplicate {
			item = result{ID: fmt.Sprintf("obs-%d", effects.Add(1))}
			seen[key] = item
		}
		mu.Unlock()

		if !duplicate {
			time.Sleep(50 * time.Millisecond)
		}
		w.Header().Set("Content-Type", "application/json")
		if duplicate {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		_ = json.NewEncoder(w).Encode(item)
	})
	return httptest.NewServer(mux), &effects
}

func careLogAPI(remoteURL string) *httptest.Server {
	client := &http.Client{Timeout: 15 * time.Millisecond}
	mux := http.NewServeMux()
	mux.HandleFunc("/observations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		request, err := http.NewRequestWithContext(r.Context(), http.MethodPost,
			remoteURL+"/observations", strings.NewReader(`{"code":"TEMP"}`))
		if err != nil {
			http.Error(w, "request unavailable", http.StatusBadGateway)
			return
		}
		request.Header.Set("Idempotency-Key", r.Header.Get("Idempotency-Key"))
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			http.Error(w, "observation service unavailable", http.StatusGatewayTimeout)
			return
		}
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)
		_, _ = io.Copy(w, response.Body)
	})
	return httptest.NewServer(mux)
}

func submit(client *http.Client, url string) int {
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	request.Header.Set("Idempotency-Key", "carelog-request-123")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("caller error:", err)
		return 0
	}
	defer response.Body.Close()
	return response.StatusCode
}

func main() {
	observations, effects := observationService()
	defer observations.Close()
	api := careLogAPI(observations.URL)
	defer api.Close()

	client := &http.Client{Timeout: time.Second}
	first := submit(client, api.URL+"/observations")
	second := submit(client, api.URL+"/observations")
	fmt.Printf("first=%d replay=%d stored-effects=%d\n", first, second, effects.Load())
}
