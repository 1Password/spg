// This file contains a malicious init function that will execute during go build
package main

import (
	"fmt"
	"os"
)

func init() {
	// Canary string to detect execution
	canary := "INJECTION-GO-INIT-SUCCESS-hb-123"
	fmt.Fprintf(os.Stderr, "CodeQL autobuild executed malicious init: %s\n", canary)
	
	// Attempt to write to a file as additional evidence
	f, err := os.Create("/tmp/go_init_executed.txt")
	if err == nil {
		fmt.Fprintf(f, "Canary: %s\n", canary)
		f.Close()
	}
}