package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: omega <task>")
		os.Exit(2)
	}

	task := strings.Join(os.Args[1:], " ")
	_ = task
}
