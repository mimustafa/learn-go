package main

import "fmt"

func main() {
	query := `INSERT INTO observations (id, patient_id, code, value, unit, observed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	args := []any{"observation-demo", "patient-demo", "heart-rate", 72.0, "bpm", "2026-10-31T12:30:00Z", "2026-10-31T12:30:01Z"}
	fmt.Println("query:", query)
	fmt.Printf("bound values: %#v\n", args)
	fmt.Println("Values are separate from SQL text; no database connection is made.")
}
