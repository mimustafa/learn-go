package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ownerOnly(next http.Handler) http.Handler {
	demoOwners := map[string]string{"patient-1": "alice"}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fixed identity is deliberately a test/demo stand-in, not authentication.
		principal := r.Header.Get("X-Demo-Principal")
		remaining := strings.TrimPrefix(r.URL.Path, "/patients/")
		patientID, _, _ := strings.Cut(remaining, "/")
		owner, exists := demoOwners[patientID]
		if principal == "" {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		if !exists || principal != owner {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	handler := ownerOnly(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("CareLog observation visible\n"))
	}))
	for _, principal := range []string{"", "alice", "bob"} {
		req := httptest.NewRequest(http.MethodGet, "/patients/patient-1/observations", nil)
		if principal != "" {
			req.Header.Set("X-Demo-Principal", principal)
		}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		fmt.Printf("principal=%q status=%d body=%s", principal, recorder.Code, strings.TrimSpace(recorder.Body.String()))
		fmt.Println()
	}
}
