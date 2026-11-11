package main

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrInvalidObservation = errors.New("invalid observation")

type ValidationError struct {
	Field string
	Why   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Why)
}

func parseObservationValue(input string) (float64, error) {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("parse observation value: %w", err)
	}
	if value < 0 {
		return 0, fmt.Errorf("%w: %w", ErrInvalidObservation,
			&ValidationError{Field: "value", Why: "must not be negative"})
	}
	return value, nil
}

func recoverDemo() {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Printf("recovered at demo boundary: %v\n", recovered)
		}
	}()
	panic("deliberate demonstration panic")
}

func main() {
	_, err := parseObservationValue("not-a-number")
	if err != nil {
		fmt.Println("parse failed:", err)
		fmt.Println("is parse error:", errors.Is(err, strconv.ErrSyntax))
	}

	_, err = parseObservationValue("-2")
	var validationErr *ValidationError
	fmt.Println("invalid observation:", errors.Is(err, ErrInvalidObservation))
	if errors.As(err, &validationErr) {
		fmt.Printf("field=%s reason=%s\n", validationErr.Field, validationErr.Why)
	}

	recoverDemo()
	fmt.Println("program continues after isolated demo")
}
