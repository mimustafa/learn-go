package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
)

type requestIDKey struct{}

func withRequestID(next http.Handler, requests *atomic.Int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" || len(id) > 64 {
			id = "generated-demo-id"
		}
		requests.Add(1)
		ctx := r.Context()
		ctx = contextWithRequestID(ctx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func main() {
	var requests atomic.Int64
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "carelog_requests_total %d\n", requests.Load())
	})
	mux.HandleFunc("/observations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, _ := r.Context().Value(requestIDKey{}).(string)
		slog.Info("observations listed", "request_id", id)
		w.Header().Set("X-Request-ID", id)
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(withRequestID(mux, &requests))
	defer server.Close()
	response, err := server.Client().Get(server.URL + "/observations")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Printf("status=%s request_id=%q requests=%d\n",
		response.Status, response.Header.Get("X-Request-ID"), requests.Load())
}
