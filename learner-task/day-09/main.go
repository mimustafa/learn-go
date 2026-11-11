package main

import (
	"errors"
	"fmt"
	"strings"
)

type Observation struct {
	Label string
	Value float64
	Unit  string
}

func (o Observation) Validate() error {
	if strings.TrimSpace(o.Label) == "" {
		return errors.New("label is required")
	}
	if strings.TrimSpace(o.Unit) == "" {
		return errors.New("unit is required")
	}
	if o.Value < 0 {
		return errors.New("value must not be negative")
	}
	return nil
}

type ObservationWriter interface {
	Save(Observation) error
}

func record(writer ObservationWriter, observation Observation) error {
	if err := observation.Validate(); err != nil {
		return fmt.Errorf("validate observation: %w", err)
	}
	return writer.Save(observation)
}

// fakeWriter records calls and injects errors without requiring a real store.
type fakeWriter struct {
	saved []Observation
	err   error
}

func (f *fakeWriter) Save(observation Observation) error {
	f.saved = append(f.saved, observation)
	return f.err
}

func main() {
	writer := &fakeWriter{}
	if err := record(writer, Observation{Label: "temperature", Value: 36.8, Unit: "°C"}); err != nil {
		fmt.Println("CareLog:", err)
		return
	}
	fmt.Printf("saved %d observation(s)\n", len(writer.saved))
}
