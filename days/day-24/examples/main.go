package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

type Store struct {
	mu   sync.Mutex
	data map[string]string
}

func (s *Store) create(id, code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = code
}

func (s *Store) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.data)
}

func handler(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/observations/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/observations/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		store.create(id, r.URL.Query().Get("code"))
		w.WriteHeader(http.StatusCreated)
	})
	return mux
}

func main() {
	store := &Store{data: make(map[string]string)}
	server := httptest.NewServer(handler(store))
	defer server.Close()

	var callers sync.WaitGroup
	for i := 0; i < 100; i++ {
		callers.Add(1)
		go func(i int) {
			defer callers.Done()
			request, _ := http.NewRequest(http.MethodPost,
				fmt.Sprintf("%s/observations/%d?code=TEMP", server.URL, i), nil)
			response, err := server.Client().Do(request)
			if err != nil {
				fmt.Println("request error:", err)
				return
			}
			response.Body.Close()
		}(i)
	}
	callers.Wait()
	fmt.Printf("created=%d\n", store.count())
}
