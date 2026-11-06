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
	AuditFields
	Label string
	Value float64
	Unit  string
}

type ObservationWriter interface {
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

func record(writer ObservationWriter, observation Observation) error {
	return writer.Save(observation)
}

func reportCount(value any) {
	reporter, ok := value.(CountReporter)
	if !ok {
		fmt.Println("count capability unavailable")
		return
	}
	fmt.Printf("stored observations: %d\n", reporter.Count())
}

func (w *MemoryWriter) Count() int {
	return len(w.Items)
}

func main() {
	writer := &MemoryWriter{}
	observation := Observation{
		AuditFields: AuditFields{RecordedAt: time.Date(2026, time.October, 31, 12, 0, 0, 0, time.UTC), RecordedBy: "demo"},
		Label:       "temperature",
		Value:       36.8,
		Unit:        "°C",
	}

	if err := record(writer, observation); err != nil {
		fmt.Println("save observation:", err)
		return
	}
	fmt.Printf("%s recorded by %s\n", observation.Label, observation.RecordedBy)
	reportCount(writer)
	reportCount(observation)
}
