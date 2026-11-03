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
## Tool definition and handler

```go
var fileReadTool = Tool{
    Name:        "file_read",
    Description: "Read the contents of a file at the given path.",
    InputSchema: map[string]any{
        "type": "object",
        "properties": map[string]any{
            "path": map[string]any{
                "type":        "string",
                "description": "Absolute file path",
            },
        },
        "required": []string{"path"},
    },
}

func runTool(name string, input json.RawMessage) (string, error) {
    switch name {
    case "file_read":
        var in struct {
            Path string `json:"path"`
        }

        if err := json.Unmarshal(input, &in); err != nil {
            return "", err
        }

        b, err := os.ReadFile(in.Path)
        return string(b), err
    }

    return "", fmt.Errorf("unknown tool: %s", name)
}
```
## The agent loop

```go
func agentLoop(client *http.Client, apiKey string, history []Message) ([]Message, error) {
    for {
        resp, err := callAPI(client, apiKey, Request{
            Model:     "claude-sonnet-5-5",
            MaxTokens: 4096,
            Messages:  history,
            Tools:     []Tool{fileReadTool},
        })

        if err != nil {
            return history, err
        }

        // Always append the assistant turn as-is
        history = append(history, Message{
            Role:    "assistant",
            Content: resp.Content,
        })

        if resp.StopReason != "tool_use" {
            return history, nil // done: print the text blocks
        }

        // Run every tool_use block, collect results in ONE user message
        var results []ContentBlock

        for _, b := range resp.Content {
            if b.Type != "tool_use" {
                continue
            }

            out, err := runTool(b.Name, b.Input)

            block := ContentBlock{
                Type:      "tool_result",
                ToolUseID: b.ID,
                Content:   out,
            }

            if err != nil {
                block.Content, block.IsError = err.Error(), true
            }

            results = append(results, block)
        }

        history = append(history, Message{
            Role:    "user",
            Content: results,
        })
    }
}
```
