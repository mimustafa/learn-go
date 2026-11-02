package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: carelog <observation-name> <value>")
		return
	}

	name, value := os.Args[1], os.Args[2]
	fmt.Printf("CareLog observation: %s = %s\n", name, value)
}
