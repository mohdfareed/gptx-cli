//go:build dev
// +build dev

package llm

import "fmt"

type DevToolArgs struct {
	Argument string `json:"argument"`
}

var DevTool = ModelTool[struct{}, DevToolArgs, *DevToolArgs]{
	Name: "dev_tool",
	Desc: "Development tool for debugging and testing.",
	Call: func(args DevToolArgs) (*DevToolArgs, error) {
		if args.Argument == "error" {
			return nil, fmt.Errorf("dev tool: %s", args.Argument)
		}
		return &args, nil
	},
}
