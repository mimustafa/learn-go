package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// This only demonstrates policy after identity is supplied; it is not authentication.
func ownerPolicy(principal, patientID string, owners map[string]string) int {
	if principal == "" {
		return http.StatusUnauthorized
	}
	if owner, ok := owners[patientID]; !ok || owner != principal {
		return http.StatusForbidden
	}
	return http.StatusOK
}

// Credential checks are inputs here; this demo does not validate real credentials.
func credentialPolicy(principal, patientID string, valid, revoked bool, owners map[string]string) int {
	if !valid || revoked {
		return http.StatusUnauthorized
	}
	return ownerPolicy(principal, patientID, owners)
}

func middleware(owners map[string]string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal := r.Header.Get("X-Demo-Principal")
		rest := strings.TrimPrefix(r.URL.Path, "/patients/")
		patientID, _, _ := strings.Cut(rest, "/")
		status := ownerPolicy(principal, patientID, owners)
		if status != http.StatusOK {
			http.Error(w, http.StatusText(status), status)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func randomOpaqueToken() (string, error) {
	// Cryptographic randomness is required for tokens; math/rand is not suitable.
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw[:]), nil
}

func main() {
	owners := map[string]string{"patient-1": "alice"}
	handler := middleware(owners, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("authorized CareLog observation access"))
	}))
	cases := []struct {
		name, principal string
		valid, revoked  bool
		want            int
	}{
		{"anonymous", "", false, false, http.StatusUnauthorized},
		{"owner", "alice", true, false, http.StatusOK},
		{"other-user", "bob", true, false, http.StatusForbidden},
		{"expired", "alice", false, false, http.StatusUnauthorized},
		{"revoked", "alice", true, true, http.StatusUnauthorized},
	}
	for _, test := range cases {
		got := credentialPolicy(test.principal, "patient-1", test.valid, test.revoked, owners)
		if got != test.want {
			panic("authorization matrix failed")
		}
		fmt.Printf("case=%s principal=%q status=%d\n", test.name, test.principal, got)
	}
	req := httptest.NewRequest(http.MethodGet, "/patients/patient-1/observations", nil)
	req.Header.Set("X-Demo-Principal", "alice")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		panic("handler policy failed")
	}
	token, err := randomOpaqueToken()
	if err != nil {
		panic(err)
	}
	fmt.Printf("example cryptographic token length=%d bytes (do not log tokens in applications)\n", len(token))
}
