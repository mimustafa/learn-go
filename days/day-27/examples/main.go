package main

import (
	"context"
	"crypto/subtle"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	start  time.Time
	used   int
	limit  int
	window time.Duration
}

func (l *Limiter) Allow(now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if now.Sub(l.start) >= l.window {
		l.start, l.used = now, 0
	}
	if l.used >= l.limit {
		return false
	}
	l.used++
	return true
}

func main() {
	wantToken := []byte("demo-token-000000")
	gotToken := []byte("demo-token-000000")
	authenticated := len(gotToken) == len(wantToken) &&
		subtle.ConstantTimeCompare(gotToken, wantToken) == 1

	limiter := &Limiter{start: time.Now(), limit: 2, window: time.Minute}
	fmt.Printf("authenticated=%t requests=%t,%t,%t\n", authenticated,
		limiter.Allow(time.Now()), limiter.Allow(time.Now()), limiter.Allow(time.Now()))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(40 * time.Millisecond)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	defer transport.CloseIdleConnections()
	client := &http.Client{Timeout: 10 * time.Millisecond, Transport: transport}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("bounded call: %v\n", err)
		return
	}
	response.Body.Close()
	fmt.Println("status:", response.Status)
}
