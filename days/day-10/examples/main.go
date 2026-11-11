package main

import "fmt"

type Number interface {
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

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, found := s[value]
	return found
}

func main() {
	fmt.Println("reading count:", Sum([]int{2, 3, 4}))
	fmt.Println("named count:", Sum([]Count{Count(5), Count(7)}))
	fmt.Println("decimal total:", Sum([]float64{36.5, 0.3}))

	labels := Set[string]{}
	labels.Add("temperature")
	fmt.Println("has temperature:", labels.Contains("temperature"))
}
