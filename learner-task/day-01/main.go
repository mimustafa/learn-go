package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func run(args []string, out, errOut io.Writer) int {
	// os.Args includes the executable name at index 0.
	if len(args) < 3 || len(args) > 4 {
		fmt.Fprintln(errOut, "usage: carelog <observation-name> <value> [unit]")
		return 2
	}

	name := strings.TrimSpace(args[1])
	if name == "" {
		fmt.Fprintln(errOut, "observation name must not be blank")
		return 2
	}

	unit := ""
	if len(args) == 4 && strings.TrimSpace(args[3]) != "" {
		unit = " " + strings.TrimSpace(args[3])
	}
	fmt.Fprintf(out, "CareLog observation: %s = %s%s\n", name, args[2], unit)
	return 0
}

func main() {
	// Keep process I/O at the CLI boundary so run can be tested independently.
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
