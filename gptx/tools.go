package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/openai/openai-go/responses"
	"github.com/urfave/cli/v3"
)

// MARK: Tools
// ============================================================================

// A model tool.
type ModelTool string
type ToolDef = responses.ToolUnionParam

// The model tool's definition.
type Tool struct {
	Name ModelTool
	Desc string
	Def  ToolDef
	Init func(Config) error
}

// The model tools manager.
type ModelTools struct {
	config Config
	tools  map[ModelTool]Tool
}

// Initialize the model tools.
func (t *ModelTools) Init() error {
	for _, tool := range t.enabledTools() {
		if err := tool.Init(t.config); err != nil {
			return fmt.Errorf("tool %s: %w", tool.Name, err)
		}
	}
	return nil
}

// MARK: CLI
// ============================================================================

func ToolsCMD(tools ModelTools) *cli.Command {
	return &cli.Command{
		Name: "tools", Usage: "show the model tools",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			toolDesc := Theme.Bold + " %s:" + Theme.Reset + " %s"

			for _, tool := range tools.enabledTools() {
				err := tool.Init(tools.config)
				if err == nil {
					fmt.Printf(
						Theme.Green+Theme.Bold+Theme.On+toolDesc,
						tool.Name, tool.Desc,
					)
				} else {
					fmt.Printf(Theme.Red+Theme.On+toolDesc, tool.Name, err)
				}
				println()
			}

			for _, tool := range tools.disabledTools() {
				fmt.Printf(
					Theme.Yellow+Theme.Bold+Theme.Off+tool.Desc,
					tool.Name, tool.Desc,
				)
				println()
			}
			for _, toolName := range tools.unknownTools() {
				fmt.Printf(
					Theme.Red+Theme.Bold+Theme.Unknown+toolDesc,
					toolName, "unknown tool",
				)
				println()
			}
			return nil
		},
	}
}

// MARK: Helpers
// ============================================================================

func (t *ModelTools) enabledTools() []Tool {
	var enabled []Tool
	for _, tool := range t.tools {
		if slices.Contains(t.config.Tools, tool.Name) {
			enabled = append(enabled, tool)
		}
	}
	return enabled
}

func (t *ModelTools) disabledTools() []Tool {
	var disabled []Tool
	for _, tool := range t.tools {
		if !slices.Contains(t.config.Tools, tool.Name) {
			disabled = append(disabled, tool)
		}
	}
	return disabled
}

func (t *ModelTools) unknownTools() []ModelTool {
	var known []ModelTool
	for _, tool := range t.tools {
		known = append(known, tool.Name)
	}

	var unknown []ModelTool
	for _, tool := range t.config.Tools {
		if !slices.Contains(known, tool) {
			unknown = append(unknown, tool)
		}
	}
	return unknown
}
