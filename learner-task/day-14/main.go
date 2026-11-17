package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	Port            int
	ShutdownTimeout time.Duration
}

func loadConfig(values map[string]string) (Config, error) {
	// Defaults apply only when an environment setting is absent.
	cfg := Config{Port: 8080, ShutdownTimeout: 10 * time.Second}
	if raw := values["PORT"]; raw != "" {
		port, err := strconv.Atoi(raw)
		if err != nil || port < 1 || port > 65535 {
			return Config{}, errors.New("PORT must be between 1 and 65535")
		}
		cfg.Port = port
	}
	if raw := values["SHUTDOWN_TIMEOUT"]; raw != "" {
		timeout, err := time.ParseDuration(raw)
		if err != nil || timeout <= 0 {
			return Config{}, errors.New("SHUTDOWN_TIMEOUT must be positive")
		}
		cfg.ShutdownTimeout = timeout
	}
	return cfg, nil
}

func serveUntil(ctx context.Context, server *http.Server, listener net.Listener, timeout time.Duration) error {
	done := make(chan error, 1)
	// Serve runs separately so this goroutine can wait for shutdown or failure.
	go func() { done <- server.Serve(listener) }()
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		// A separate finite deadline bounds how long graceful shutdown can wait.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		err := <-done
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func configuredServer(handler http.Handler, address string) *http.Server {
	return &http.Server{
		Addr: address, Handler: handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func signalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func cancelableMutation(w http.ResponseWriter, r *http.Request, mutate func()) {
	select {
	case <-r.Context().Done():
		return
	case <-time.After(time.Second):
		mutate()
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	cfg, err := loadConfig(map[string]string{"PORT": "8082", "SHUTDOWN_TIMEOUT": "3s"})
	if err != nil {
		panic(err)
	}
	if _, err := loadConfig(map[string]string{"PORT": "70000"}); err == nil {
		panic("invalid port accepted")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/v1/observations", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("[]")) })
	_ = configuredServer(mux, ":8082")
	for _, path := range []string{"/healthz", "/readyz", "/v1/observations"} {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		fmt.Printf("%s=%d ", path, recorder.Code)
	}
	mutations := 0
	requestCtx, cancelRequest := context.WithCancel(context.Background())
	cancelRequest()
	recorder := httptest.NewRecorder()
	cancelableMutation(recorder, httptest.NewRequest(http.MethodPost, "/v1/observations", nil).WithContext(requestCtx),
		func() { mutations++ })
	if mutations != 0 {
		panic("canceled request mutated application state")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	<-ctx.Done()
	fmt.Printf("config=%+v cancellation=%v\n", cfg, ctx.Err())
}
