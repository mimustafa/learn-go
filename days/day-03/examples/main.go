package main

import "fmt"

func main() {
	label := "café"
	fmt.Printf("%q has %d UTF-8 bytes\n", label, len(label))

	for offset, codePoint := range label {
		fmt.Printf("byte offset %d: %U\n", offset, codePoint)
	}

	readings := []int{118, 121, 119}
	total := 0
	for _, reading := range readings {
		total += reading
	}
	fmt.Printf("readings=%v total=%d\n", readings, total)

	counts := map[string]int{"blood pressure": len(readings)}
	count, found := counts["blood pressure"]
	fmt.Printf("blood pressure count=%d found=%t\n", count, found)

	_, found = counts["temperature"]
	fmt.Printf("temperature present=%t\n", found)
}
