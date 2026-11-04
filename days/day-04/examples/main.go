package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseValue(input string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return 0, fmt.Errorf("parse observation value %q: %w", input, err)
	}
	return value, nil
}

func summarize(values ...float64) (count int, total float64) {
	for _, value := range values {
		count++
		total += value
	}
	return count, total
}

func main() {
	defer fmt.Println("finished CareLog example")

	labelFor := func(encounter string) func(string) string {
		return func(name string) string {
			return encounter + ": " + name
		}
	}

	value, err := parseValue(" 118.5 ")
	if err != nil {
		fmt.Println(err)
		return
	}

	count, total := summarize(value, 120, 117.5)
	fmt.Printf("%s = %.1f; count=%d total=%.1f\n",
		labelFor("visit-1")("blood pressure"), value, count, total)
}
