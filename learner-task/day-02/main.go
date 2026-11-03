package main

import "fmt"

const demonstrationLimit float64 = 100

func classify(value float64) string {
	if value < 0 {
		return "invalid negative value"
	}

	switch {
	case value < demonstrationLimit:
		return "below demonstration limit"
	case value == demonstrationLimit:
		return "at demonstration limit"
	default:
		return "above demonstration limit"
	}
}

func main() {
	values := []float64{72, 100, 124, -1}
	withinLimit := 0

	for _, value := range values {
		fmt.Printf("%.0f: %s\n", value, classify(value))
		if value >= 0 && value <= demonstrationLimit {
			withinLimit++
		}
	}

	fmt.Printf("%d of %d values are non-negative and at or below the demonstration limit\n", withinLimit, len(values))
}
