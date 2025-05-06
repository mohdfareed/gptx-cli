package llm

// ModelConfig is the model's generation configuration.
type ModelConfig struct {
	SysPrompt string `json:"prompt"`
	Tools     []Tool `json:"tools"`
}

// Model is the LLM model interface.
type Model struct {
	Started    *Event[struct{}]             // start
	Generated  *Event[string]               // stream chunk
	ToolCall   *Event[ToolCall[any]]        // tool, args
	ToolResult *Event[ToolResult[any, any]] // tool, args, result
	Finished   *Event[Msg]                  // end
	Error      *Event[error]                // error
	config     ModelConfig
}

// Usage is the model's usage statistics.
type Usage struct {
	InputTokens int64 `json:"input_tokens"`
	GenTokens   int64 `json:"generated_tokens"`
}

// TotalUsage returns the total aggregated usage of the messages.
func TotalUsage(msgs []Msg) Usage {
	usage := Usage{}
	for _, msg := range msgs {
		usage.InputTokens += msg.UsageData().InputTokens
		usage.GenTokens += msg.UsageData().GenTokens
	}
	return usage
}
