package llm

import (
	"fmt"
	"slices"
)

// Tool is a model tool's name.
type Tool string

// ModelTool is a model tool's definition.
type ModelTool[Config any, Params any, Result any] struct {
	Name Tool
	Desc string
	Call func(Params) (Result, error)
	Init func() error
	// config Config
}

// ModelTools is the tools used by the model.
type ModelTools map[Tool]ModelTool[any, any, any]

// ToolCall is a model tool's call and its results.
type ToolCall[Params any, Results any] struct {
	Tool   Tool    `json:"tool"`
	Args   Params  `json:"args"`
	Result Results `json:"results"`
}

// MARK: Init
// ============================================================================

// Create new model tools.
func CreateTools(tools ModelTools, model *Model) ModelTools {
	tools = tools.EnabledTools(model.config)
	model.Started.Subscribe(func(struct{}) {
		for _, tool := range tools.EnabledTools(model.config) {
			if err := tool.Init(); err != nil {
				model.Error.publish(err)
			}
		}
	})

	model.ToolCall.Subscribe(func(args P2[Tool, any]) {
		result, err := tools.call(args.A, args.B)
		if err != nil {
			model.Error.publish(err)
			return
		}
		model.ToolResult.publish(ToolCall[any, any]{
			Tool:   args.A,
			Args:   args.B,
			Result: result,
		})
	})
	return tools
}

func (t *ModelTools) call(tool Tool, args any) (any, error) {
	for _, t := range *t {
		if t.Name == tool {
			return t.Call(args)
		}
	}
	return nil, fmt.Errorf("tool %s: not found", tool)
}

// MARK: Helpers
// ============================================================================

func (t *ModelTools) EnabledTools(config ModelConfig) ModelTools {
	enabled := make(ModelTools)
	for _, tool := range *t {
		if slices.Contains(config.Tools, tool.Name) {
			enabled[tool.Name] = tool
		}
	}
	return enabled
}

func (t *ModelTools) DisabledTools(config ModelConfig) ModelTools {
	disabled := make(ModelTools)
	for _, tool := range *t {
		if !slices.Contains(config.Tools, tool.Name) {
			disabled[tool.Name] = tool
		}
	}
	return disabled
}

func (t *ModelTools) UnknownTools(config ModelConfig) []Tool {
	var known []Tool
	for _, tool := range *t {
		known = append(known, tool.Name)
	}

	var unknown []Tool
	for _, tool := range config.Tools {
		if !slices.Contains(known, tool) {
			unknown = append(unknown, tool)
		}
	}
	return unknown
}
