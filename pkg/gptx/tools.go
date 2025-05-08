package gptx

import (
	"fmt"
	"slices"
)

// Tool is a model tool's name.
type Tool string

// ModelTool is a model tool's definition.
type ModelTool[Args any] struct {
	Name   Tool
	Desc   string
	Call   func(Args) error
	Result chan any
}

// ModelTools is the tools used by the model.
type ModelTools map[Tool]AnyModelTool

// AnyModelTool is a model tool's with any parameters and result.
type AnyModelTool = ModelTool[any, any]

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

// // MARK: Tools
// // ============================================================================

// // A model tool.
// type Tool string
// type ToolDef = responses.ToolUnionParam

// // The model tool's definition.
// type ModelTool struct {
// 	Name Tool
// 	Desc string
// 	Def  ToolDef
// 	Init func(Config) error
// }

// // The model tools manager.
// type ModelTools struct {
// 	config Config
// 	tools  map[Tool]ModelTool
// }

// // Initialize the model tools.
// func (t *ModelTools) Init() error {
// 	for _, tool := range t.enabledTools() {
// 		if err := tool.Init(t.config); err != nil {
// 			return fmt.Errorf("tool %s: %w", tool.Name, err)
// 		}
// 	}
// 	return nil
// }

// // MARK: CLI
// // ============================================================================

// func ToolsCMD(tools ModelTools) *cli.Command {
// 	return &cli.Command{
// 		Name: "tools", Usage: "show the model tools",
// 		Action: func(ctx context.Context, cmd *cli.Command) error {
// 			toolDesc := Theme.Bold + " %s:" + Theme.Reset + " %s"

// 			for _, tool := range tools.enabledTools() {
// 				err := tool.Init(tools.config)
// 				if err == nil {
// 					fmt.Printf(
// 						Theme.Green+Theme.Bold+Theme.On+toolDesc,
// 						tool.Name, tool.Desc,
// 					)
// 				} else {
// 					fmt.Printf(Theme.Red+Theme.On+toolDesc, tool.Name, err)
// 				}
// 				println()
// 			}

// 			for _, tool := range tools.disabledTools() {
// 				fmt.Printf(
// 					Theme.Yellow+Theme.Bold+Theme.Off+tool.Desc,
// 					tool.Name, tool.Desc,
// 				)
// 				println()
// 			}
// 			for _, toolName := range tools.unknownTools() {
// 				fmt.Printf(
// 					Theme.Red+Theme.Bold+Theme.Unknown+toolDesc,
// 					toolName, "unknown tool",
// 				)
// 				println()
// 			}
// 			return nil
// 		},
// 	}
// }

// // MARK: Helpers
// // ============================================================================

// func (t *ModelTools) enabledTools() []ModelTool {
// 	var enabled []ModelTool
// 	for _, tool := range t.tools {
// 		if slices.Contains(t.config.Tools, tool.Name) {
// 			enabled = append(enabled, tool)
// 		}
// 	}
// 	return enabled
// }

// func (t *ModelTools) disabledTools() []ModelTool {
// 	var disabled []ModelTool
// 	for _, tool := range t.tools {
// 		if !slices.Contains(t.config.Tools, tool.Name) {
// 			disabled = append(disabled, tool)
// 		}
// 	}
// 	return disabled
// }

// func (t *ModelTools) unknownTools() []Tool {
// 	var known []Tool
// 	for _, tool := range t.tools {
// 		known = append(known, tool.Name)
// 	}

// 	var unknown []Tool
// 	for _, tool := range t.config.Tools {
// 		if !slices.Contains(known, tool) {
// 			unknown = append(unknown, tool)
// 		}
// 	}
// 	return unknown
// }
