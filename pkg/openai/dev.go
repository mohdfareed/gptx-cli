//go:build dev
// +build dev

package openai

import "fmt"

type DevToolArgs struct {
	Argument string `json:"argument"`
}

var DevOpenAI = ModelTool[DevToolArgs]{
	Name: "dev_tool",
	Desc: "Development tool for debugging and testing.",
	Call: func(args DevToolArgs) (any, error) {
		if args.Argument == "error" {
			return nil, fmt.Errorf("dev tool: %s", args.Argument)
		}
		return args, nil
	},
}
