## Types

```go
package main

import "encoding/json"

// One content block. "type" decides which fields are used.
type ContentBlock struct {
    Type string `json:"type"` // "text" | "tool_use" | "tool_result"

    // text
    Text string `json:"text,omitempty"`

    // tool_use (assistant -> you)
    ID    string          `json:"id,omitempty"`
    Name  string          `json:"name,omitempty"`
    Input json.RawMessage `json:"input,omitempty"` // decode per tool

    // tool_result (you -> assistant)
    ToolUseID string `json:"tool_use_id,omitempty"`
    Content   string `json:"content,omitempty"`
    IsError   bool   `json:"is_error,omitempty"`
}

type Message struct {
    Role    string         `json:"role"` // "user" | "assistant"
    Content []ContentBlock `json:"content"`
}

type Tool struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    InputSchema map[string]any `json:"input_schema"`
}

type Request struct {
    Model     string    `json:"model"`
    MaxTokens int       `json:"max_tokens"`
    System    string    `json:"system,omitempty"`
    Messages  []Message `json:"messages"`
    Tools     []Tool    `json:"tools,omitempty"`
}

// What the API returns
type Response struct {
    ID         string         `json:"id"`
    Role       string         `json:"role"`
    Content    []ContentBlock `json:"content"`
    StopReason string         `json:"stop_reason"` // "end_turn" | "tool_use" | "max_tokens"
}
```
