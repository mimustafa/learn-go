package main

import (
	"fmt"
	"strings"
	"time"
)

// A defined type gives IDs their own type instead of treating them as any string.
type ObservationID string

type Observation struct {
	ID       ObservationID
	Label    string
	Value    float64
	Unit     string
	Recorded time.Time
}

func (o Observation) Validate() error {
	if strings.TrimSpace(string(o.ID)) == "" {
		return fmt.Errorf("observation ID is required")
	}
	if strings.TrimSpace(o.Label) == "" {
		return fmt.Errorf("observation label is required")
	}
	if strings.TrimSpace(o.Unit) == "" {
		return fmt.Errorf("observation unit is required")
	}
	return nil
}

func (o Observation) Display() string {
	return fmt.Sprintf("%s: %.1f %s", o.Label, o.Value, o.Unit)
}

func (o Observation) NormalizeLabel() Observation {
	// A value receiver changes and returns a copy, leaving the original untouched.
	o.Label = strings.TrimSpace(o.Label)
	return o
}

func main() {
	observation := Observation{
		ID:       ObservationID("obs-1"),
		Label:    " temperature ",
		Value:    36.8,
		Unit:     "°C",
		Recorded: time.Date(2026, time.October, 31, 12, 0, 0, 0, time.UTC),
	}
	normalized := observation.NormalizeLabel()
	if err := normalized.Validate(); err != nil {
		fmt.Println("invalid observation:", err)
		return
	}
	fmt.Println(normalized.Display())
	fmt.Printf("original label remains %q\n", observation.Label)
}
