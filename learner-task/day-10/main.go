package main

import (
	"fmt"
	"strings"
)

type Number interface {
	// ~ also admits named types whose underlying type is int or float64.
	~int | ~float64
}

type Count int

func Sum[T Number](values []T) T {
	var total T
	for _, value := range values {
		total += value
	}
	return total
}

// An empty struct records membership without storing an extra value per key.
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, found := s[value]
	return found
}

type Observation struct {
	Label string
	Value float64
	Unit  string
}

func formatObservation(observation Observation) string {
	return fmt.Sprintf("%s: %.1f %s", strings.TrimSpace(observation.Label), observation.Value, observation.Unit)
}

func main() {
	fmt.Println("integer total:", Sum([]int{2, 3, 4}))
	fmt.Println("named integer total:", Sum([]Count{Count(5), Count(7)}))
	fmt.Println("decimal total:", Sum([]float64{1.25, 2.5}))

	labels := Set[string]{}
	labels.Add("temperature")
	fmt.Println("contains temperature:", labels.Contains("temperature"))
	fmt.Println(formatObservation(Observation{Label: " temperature ", Value: 36.8, Unit: "°C"}))
}
