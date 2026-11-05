package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseValue(input string) (float64, error) {
	if strings.TrimSpace(input) == "" {
		return 0, errors.New("observation value must not be blank")
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		// %w keeps the parse error available to errors.Is and errors.As.
		return 0, fmt.Errorf("parse observation value %q: %w", input, err)
	}
	return value, nil
}

func summarize(values ...float64) (count int, total, average float64) {
	// A variadic parameter is received as a slice; no values means zero results.
	for _, value := range values {
		count++
		total += value
	}
	if count > 0 {
		average = total / float64(count)
	}
	return count, total, average
}

func main() {
	// Deferred calls run in last-in, first-out order.
	defer fmt.Println("outer cleanup")
	defer fmt.Println("inner cleanup")

	labelFor := func(encounter string) func(string) string {
		// The returned function closes over the encounter label.
		return func(name string) string { return encounter + ": " + name }
	}
	value, err := parseValue(" 118.5 ")
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	count, total, average := summarize(value, 120, 117.5)
	fmt.Printf("%s; count=%d total=%.1f average=%.2f\n",
		labelFor("visit-2")("blood pressure"), count, total, average)
	emptyCount, emptyTotal, emptyAverage := summarize()
	fmt.Printf("empty summary: count=%d total=%.1f average=%.1f\n",
		emptyCount, emptyTotal, emptyAverage)
}
