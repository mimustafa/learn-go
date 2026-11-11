package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// Sentinels let callers check error categories with errors.Is.
	ErrInvalidObservation = errors.New("invalid observation")
	ErrMissingLabel       = errors.New("missing observation label")
)

type ValidationError struct {
	Field string
	Why   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Why)
}

func parseObservationValue(input string) (float64, error) {
	if strings.TrimSpace(input) == "" {
		return 0, fmt.Errorf("parse observation value: input is blank")
	}
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("parse observation value: %w", err)
	}
	if value < 0 {
		validationErr := &ValidationError{Field: "value", Why: "must not be negative"}
		// Wrapping both errors preserves the category and its structured details.
		return 0, fmt.Errorf("%w: %w", ErrInvalidObservation, validationErr)
	}
	return value, nil
}

func validateLabel(label string) error {
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("%w", ErrMissingLabel)
	}
	return nil
}

func recoverDemo() {
	defer func() {
		// recover only catches a panic in this goroutine's deferred call.
		if recovered := recover(); recovered != nil {
			fmt.Printf("recovered at demo boundary: %v\n", recovered)
		}
	}()
	panic("deliberate demonstration panic")
}

func main() {
	_, err := parseObservationValue("not-a-number")
	if err != nil {
		fmt.Println("cannot parse observation:", err)
		fmt.Println("is syntax error:", errors.Is(err, strconv.ErrSyntax))
	}
	if err := validateLabel(" "); err != nil {
		fmt.Println("is missing label:", errors.Is(err, ErrMissingLabel))
	}
	_, err = parseObservationValue("-2")
	var validationErr *ValidationError
	fmt.Println("is invalid observation:", errors.Is(err, ErrInvalidObservation))
	if errors.As(err, &validationErr) {
		fmt.Printf("field=%s reason=%s\n", validationErr.Field, validationErr.Why)
	}
	recoverDemo()
	fmt.Println("program continues after isolated recovery")
}
