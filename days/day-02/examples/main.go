package main

import "fmt"

const demonstrationLimit float64 = 100

func main() {
	values := []float64{72, 100, 124}
	withinLimit := 0

	for _, value := range values {
		if value < 0 {
			fmt.Printf("%.0f: invalid negative value\n", value)
			continue
		}

		switch {
		case value < demonstrationLimit:
			fmt.Printf("%.0f: below demonstration limit\n", value)
		case value == demonstrationLimit:
			fmt.Printf("%.0f: at demonstration limit\n", value)
		default:
			fmt.Printf("%.0f: above demonstration limit\n", value)
		}
		if value <= demonstrationLimit {
			withinLimit++
		}
	}

	fmt.Printf("%d of %d values are at or below the demonstration limit\n", withinLimit, len(values))
}
