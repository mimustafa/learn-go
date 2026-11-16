package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("live\n"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("/v1/observations", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"observations":[]}`))
	})

	for _, path := range []string{"/healthz", "/readyz", "/v1/observations"} {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		fmt.Printf("%s -> %d %s", path, recorder.Code, recorder.Body.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	<-ctx.Done()
	fmt.Printf("request budget ended: %v\n", ctx.Err())
}
