package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: omega <task>")
		os.Exit(2)
	}
	task := strings.Join(os.Args[1:], " ")

	key := os.Getenv("OPENROUTER_API_KEY")
	model := os.Getenv("OPENROUTER_MODEL")
	if key == "" || model == "" {
		fmt.Fprintln(os.Stderr, "error: OPENROUTER_API_KEY or OPENROUTER_MODEL is not set")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// temporary: silence "declared and not used" errors
	_ = task
	_ = key
	_ = model
	_ = ctx
	system := "You are a coding agent with access to tools. When the user mentions a file, you MUST call the read_file tool first. Example: if the user says \"read go.mod\", call read_file with {\"file_path\": \"go.mod\"}. Then answer based only on the tool result."
	_ = system // temporary
}
