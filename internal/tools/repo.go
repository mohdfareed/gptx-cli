package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
)

const RepoToolName = "repo"
const RepoToolDescription = `Interact with a codebase.
Use this to list directories and read files.
`

var repoTool = ToolDef{
	Name: RepoToolName,
	Desc: RepoToolDescription,
	Params: map[string]any{
		"path":     []string{},
		"contents": false,
	},
	Handler: repoHandler,
}

func RepoTool(config cfg.Config) ToolDef {
	tool := repoTool
	tool.Params["path"] = config.Repo
	return tool
}

func repoHandler(ctx context.Context, params map[string]any) (string, error) {
	base := "." // or configurable root
	target := filepath.Join(base, params["paths"].(string))

	if params["contents"].(bool) {
		data, err := os.ReadFile(target)
		if err != nil {
			return "", fmt.Errorf("read: %w", err)
		}
		resp := map[string]any{
			"path":    params["path"],
			"content": string(data),
		}

		// Marshal the response to JSON
		j, err := json.Marshal(resp)
		if err != nil {
			return "", fmt.Errorf("json: %w", err)
		}
		return string(j), nil
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		return "", fmt.Errorf("list: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		names = append(names, name)
	}
	resp := map[string]any{
		"path":    params["path"],
		"entries": names,
	}

	// Marshal the response to JSON
	j, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("json: %w", err)
	}
	return string(j), nil

}
