package main

import (
	"fmt"
	"strconv"
	"strings"
)

func observationSummary(label, rawValue string) (string, error) {
	label = strings.TrimSpace(label)
	value, err := strconv.ParseFloat(strings.TrimSpace(rawValue), 64)
	if err != nil {
		return "", fmt.Errorf("parse value: %w", err)
	}
	if label == "" {
		return "", fmt.Errorf("label is required")
	}
	return fmt.Sprintf("%s = %.1f", label, value), nil
}

func main() {
	summary, err := observationSummary(" temperature ", "36.8")
	if err != nil {
		fmt.Println("CareLog:", err)
		return
	}
	fmt.Println("CareLog observation:", summary)
}
