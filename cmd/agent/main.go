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
	key := os.Getenv("OPENROUTER_API_KEY")
	model := os.Getenv("OPENROUTER_MODEL")
	if key == "" || model == "" {
		fmt.Fprintln(os.Stderr, "error: OPENROUTER_API_KEY or OPENROUTER_MODEL is not set")
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	_ = ctx // temporary