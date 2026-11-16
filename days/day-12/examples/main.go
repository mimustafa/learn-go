package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type patientRequest struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type patientResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

func main() {
	var patients []patientResponse
	mux := http.NewServeMux()
	mux.HandleFunc("/patients", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(patients)
			return
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var input patientRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		patient := patientResponse{ID: fmt.Sprintf("patient-%d", len(patients)+1), Name: input.Name, BirthDate: input.BirthDate}
		patients = append(patients, patient)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(patient)
	})

	post := httptest.NewRequest(http.MethodPost, "/patients", strings.NewReader(`{"name":"Ada Example","birth_date":"1815-12-10"}`))
	created := httptest.NewRecorder()
	mux.ServeHTTP(created, post)
	listed := httptest.NewRecorder()
	mux.ServeHTTP(listed, httptest.NewRequest(http.MethodGet, "/patients", nil))
	fmt.Printf("POST status %d; patients: %s", created.Code, listed.Body.String())
}
