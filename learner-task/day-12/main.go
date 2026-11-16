package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type patient struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Secret    string
}

type patientResponse struct {
	// The response DTO controls which fields are exposed to API clients.
	ID        string `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type patientInput struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type patientService struct{ patients map[string]patient }

func responseFor(p patient) patientResponse {
	return patientResponse{ID: p.ID, Name: p.Name, BirthDate: p.BirthDate}
}

func (s *patientService) handler() http.Handler {
	mux := http.NewServeMux()
	// Method-aware patterns keep GET and POST contracts separate.
	mux.HandleFunc("GET /patients", func(w http.ResponseWriter, _ *http.Request) {
		list := make([]patientResponse, 0, len(s.patients))
		for _, item := range s.patients {
			list = append(list, responseFor(item))
		}
		writeJSON(w, http.StatusOK, list)
	})
	mux.HandleFunc("POST /patients", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var input patientInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		id := fmt.Sprintf("patient-%d", len(s.patients)+1)
		item := patient{ID: id, Name: input.Name, BirthDate: input.BirthDate, Secret: "internal"}
		s.patients[id] = item
		writeJSON(w, http.StatusCreated, responseFor(item))
	})
	mux.HandleFunc("GET /patients/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, ok := s.patients[r.PathValue("id")]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, responseFor(item))
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	// Set headers before writing the status or response body.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func main() {
	service := &patientService{patients: make(map[string]patient)}
	handler := service.handler()
	bad := httptest.NewRecorder()
	handler.ServeHTTP(bad, httptest.NewRequest(http.MethodPost, "/patients", strings.NewReader("{")))
	if bad.Code != http.StatusBadRequest || len(service.patients) != 0 {
		panic("malformed request mutated the service")
	}
	created := httptest.NewRecorder()
	handler.ServeHTTP(created, httptest.NewRequest(http.MethodPost, "/patients",
		strings.NewReader(`{"name":"Ada Example","birth_date":"1815-12-10"}`)))
	if created.Code != http.StatusCreated || strings.Contains(created.Body.String(), "internal") {
		panic("response DTO/status incorrect")
	}
	found := httptest.NewRecorder()
	handler.ServeHTTP(found, httptest.NewRequest(http.MethodGet, "/patients/patient-1", nil))
	missing := httptest.NewRecorder()
	handler.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/patients/missing", nil))
	fmt.Printf("POST=%d GET=%d missing=%d response=%s", created.Code, found.Code, missing.Code, created.Body.String())
}
