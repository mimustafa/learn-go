package main

import "fmt"

func countRunes(value string) int {
	count := 0
	for range value {
		count++
	}
	return count
}

func main() {
	label := "café"
	fmt.Printf("%q has %d UTF-8 bytes and %d runes\n", label, len(label), countRunes(label))
	fmt.Printf("bytes=%v runes=%U\n", []byte(label), []rune(label))
	for offset, codePoint := range label {
		fmt.Printf("byte offset %d: %U\n", offset, codePoint)
	}

	readings := []int{118, 121, 119}
	readings = append(readings, 120)
	total := 0
	for _, reading := range readings {
		total += reading
	}
	fmt.Printf("readings=%v total=%d\n", readings, total)

	counts := map[string]int{"blood pressure": len(readings), "temperature": 0}
	for _, key := range []string{"blood pressure", "temperature", "oxygen"} {
		count, found := counts[key]
		fmt.Printf("%s count=%d found=%t\n", key, count, found)
	}
}
