package main

import (
	"fmt"
	"strings"
	"time"
)

type Observation struct {
	Label    string
	Value    float64
	Unit     string
	Recorded time.Time
}

func (o Observation) Validate() error {
	if strings.TrimSpace(o.Label) == "" {
		return fmt.Errorf("observation label is required")
	}
	if strings.TrimSpace(o.Unit) == "" {
		return fmt.Errorf("observation unit is required")
	}
	return nil
}

func (o *Observation) SetUnit(unit string) {
	o.Unit = strings.TrimSpace(unit)
}

func (o Observation) Display() string {
	return fmt.Sprintf("%s: %.1f %s", o.Label, o.Value, o.Unit)
}

func main() {
	observation := Observation{
		Label:    "temperature",
		Value:    36.8,
		Recorded: time.Date(2026, time.October, 31, 12, 0, 0, 0, time.UTC),
	}
	observation.SetUnit(" °C ")

	if err := observation.Validate(); err != nil {
		fmt.Println("invalid observation:", err)
		return
	}
	fmt.Println(observation.Display())
}
