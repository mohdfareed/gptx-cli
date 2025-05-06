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

// MARK: Init
// ============================================================================

// Create new model tools.
func CreateTools(tools ModelTools, model *Model) ModelTools {
	tools = tools.enabledTools(model.config)
	model.Started.Subscribe(func(struct{}) {
		for _, tool := range tools.enabledTools(model.config) {
			if err := tool.Init(); err != nil {
				model.Error.publish(err)
			}
		}
	})

	model.ToolCall.Subscribe(func(args ToolCall[any]) {
		result, err := tools.call(args.Tool, args.Args)
		if err != nil {
			model.Error.publish(err)
			return
		}
		model.ToolResult.publish(ToolResult[any, any]{
			ID:     args.ID,
			Result: result,
		})
	})
	return tools
}

// MARK: Helpers
// ============================================================================

func (t *ModelTools) call(tool Tool, args any) (any, error) {
	for _, t := range *t {
		if t.Name == tool {
			return t.Call(args)
		}
	}
	return nil, fmt.Errorf("tool %s: not found", tool)
}

func (t *ModelTools) enabledTools(config ModelConfig) ModelTools {
	enabled := make(ModelTools)
	for _, tool := range *t {
		if slices.Contains(config.Tools, tool.Name) {
			enabled[tool.Name] = tool
		}
	}
	return enabled
}
