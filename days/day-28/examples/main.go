package main

import (
	"errors"
	"fmt"
)

type Patient struct {
	ID string
}

type PatientReader interface {
	Exists(id string) bool
}

type Observation struct {
	PatientID string
	Code      string
	Value     float64
}

type ObservationModule struct {
	patients PatientReader
	items    []Observation
}

func (m *ObservationModule) Record(item Observation) error {
	if !m.patients.Exists(item.PatientID) {
		return errors.New("patient not found")
	}
	if item.Code == "" {
		return errors.New("code is required")
	}
	m.items = append(m.items, item)
	return nil
}

type inMemoryPatients map[string]Patient

func (p inMemoryPatients) Exists(id string) bool {
	_, ok := p[id]
	return ok
}

func main() {
	patients := inMemoryPatients{"patient-1": {ID: "patient-1"}}
	observations := ObservationModule{patients: patients}
	err := observations.Record(Observation{
		PatientID: "patient-1", Code: "TEMP", Value: 37.2,
	})
	fmt.Printf("recorded=%t error=%v\n", err == nil, err)
}
