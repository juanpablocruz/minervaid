package main

import (
	"fmt"
	"os"

	"github.com/juanpablocruz/minervaid/cmd/ego/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// Print the error to stderr and exit with code 1
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
