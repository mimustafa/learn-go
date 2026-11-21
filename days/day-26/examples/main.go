package main

import (
	"bytes"
	"fmt"
	"runtime/pprof"
	"time"
)

func cpuWork() uint64 {
	var total uint64
	for i := uint64(1); i < 2_000_000; i++ {
		total += (i * 31) % 97
	}
	return total
}

func main() {
	var cpu bytes.Buffer
	if err := pprof.StartCPUProfile(&cpu); err != nil {
		panic(err)
	}
	value := cpuWork()
	time.Sleep(100 * time.Millisecond)
	pprof.StopCPUProfile()

	var heap bytes.Buffer
	if err := pprof.Lookup("heap").WriteTo(&heap, 0); err != nil {
		panic(err)
	}
	fmt.Printf("work=%d cpu-profile-bytes=%d heap-profile-bytes=%d\n",
		value, cpu.Len(), heap.Len())
}
