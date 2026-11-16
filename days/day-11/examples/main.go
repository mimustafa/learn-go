package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Observation struct {
	ID         string  `json:"id"`
	PatientID  string  `json:"patient_id"`
	Code       string  `json:"code"`
	Value      float64 `json:"value"`
	Unit       string  `json:"unit"`
	ObservedAt string  `json:"observed_at"`
	CreatedAt  string  `json:"created_at"`
}

func main() {
	path := "carelog-demo.json"
	defer os.Remove(path)

	observations := []Observation{
		{ID: "observation-1", PatientID: "patient-demo", Code: "heart-rate", Value: 72, Unit: "bpm", ObservedAt: "2026-10-31T12:30:00Z", CreatedAt: "2026-10-31T12:30:01Z"},
		{ID: "observation-2", PatientID: "patient-demo", Code: "sleep-hours", Value: 8, Unit: "h", ObservedAt: "2026-10-31T07:00:00Z", CreatedAt: "2026-10-31T12:30:01Z"},
	}
	data, err := json.MarshalIndent(observations, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		panic(err)
	}

	stored, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var restored []Observation
	if err := json.Unmarshal(stored, &restored); err != nil {
		panic(err)
	}
	fmt.Printf("restored %d observations; first: %s = %.0f %s\n", len(restored), restored[0].Code, restored[0].Value, restored[0].Unit)
}
