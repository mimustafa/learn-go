package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Observation struct {
	Label string
	Value float64
}

func Validate(observation Observation) error {
	if normalizeLabel(observation.Label) == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}

// Package-private helper: callers use the exported validation behavior instead.
func normalizeLabel(label string) string {
	return strings.TrimSpace(label)
}

func summary(label, rawValue string) (string, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(rawValue), 64)
	if err != nil {
		return "", fmt.Errorf("parse observation value: %w", err)
	}
	observation := Observation{Label: normalizeLabel(label), Value: value}
	if err := Validate(observation); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s = %.1f", observation.Label, observation.Value), nil
}

func main() {
	value, err := summary(" temperature ", "36.8")
	if err != nil {
		fmt.Println("CareLog:", err)
		return
	}
	fmt.Println("CareLog observation:", value)
}
