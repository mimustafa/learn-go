package main

import (
	"fmt"
	"runtime"
	"strings"
)

func concatenate(parts []string) string {
	var builder strings.Builder
	size := 0
	for _, part := range parts {
		size += len(part)
	}
	builder.Grow(size)
	for _, part := range parts {
		builder.WriteString(part)
	}
	return builder.String()
}

func main() {
	parts := make([]string, 1000)
	for i := range parts {
		parts[i] = "observation;"
	}

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	result := concatenate(parts)
	runtime.ReadMemStats(&after)

	values := make([]int, 4, 8)
	values = append(values, 5)
	fmt.Printf("bytes=%d slice-len=%d cap=%d\n", len(result), len(values), cap(values))
	fmt.Printf("total-alloc-delta=%d bytes (process-wide observation)\n",
		after.TotalAlloc-before.TotalAlloc)
}
