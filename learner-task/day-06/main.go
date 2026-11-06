package main

import (
	"fmt"
	"time"
)

type AuditFields struct {
	RecordedAt time.Time
	RecordedBy string
}

type Observation struct {
	// Embedding promotes these fields, but does not create an inheritance hierarchy.
	AuditFields
	Label string
	Value float64
	Unit  string
}

type ObservationWriter interface {
	// Implementations satisfy this interface by providing Save; no declaration is needed.
	Save(Observation) error
}

type MemoryWriter struct {
	Items []Observation
}

func (w *MemoryWriter) Save(observation Observation) error {
	w.Items = append(w.Items, observation)
	return nil
}

type CountReporter interface {
	Count() int
}

func (w *MemoryWriter) Count() int {
	return len(w.Items)
}

func record(writer ObservationWriter, observation Observation) error {
	return writer.Save(observation)
}

func reportCapabilities(value any) {
	// A type switch checks optional behavior without risking a failed assertion panic.
	switch capability := value.(type) {
	case CountReporter:
		fmt.Printf("count capability: %d\n", capability.Count())
	case fmt.Stringer:
		fmt.Println("string capability:", capability.String())
	default:
		fmt.Println("optional capability unavailable")
	}
}

func main() {
	writer := &MemoryWriter{}
	observation := Observation{
		AuditFields: AuditFields{
			RecordedAt: time.Date(2026, time.October, 31, 12, 0, 0, 0, time.UTC),
			RecordedBy: "demo",
		},
		Label: "temperature",
		Value: 36.8,
		Unit:  "°C",
	}
	if err := record(writer, observation); err != nil {
		fmt.Println("save observation:", err)
		return
	}
	fmt.Printf("%s recorded by %s\n", observation.Label, observation.RecordedBy)
	reportCapabilities(writer)
	reportCapabilities(observation)
}
