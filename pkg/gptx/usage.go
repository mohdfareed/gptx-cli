package gptx

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohdfareed/gptx-cli/pkg/openai"
)

// UsagePath is the path to the usage file.
var UsagePath = appDir + "/usage.json"

// GetUsage reads the current usage from the usage file and returns it.
// If the file does not exist, it returns a zero MsgUsage.
func GetUsage() (openai.MsgUsage, error) {
	var u openai.MsgUsage
	data, err := os.ReadFile(UsagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return openai.MsgUsage{}, nil
		}
		return u, fmt.Errorf("read usage: %w", err)
	}

	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("unmarshal usage: %w", err)
	}
	return u, nil
}

// AddUsage loads the existing usage, adds the provided MsgUsage values,
// and writes the updated totals back to the usage file.
func AddUsage(usage openai.MsgUsage) error {
	curr, err := GetUsage()
	if err != nil {
		return err
	}
	total := openai.TotalUsage([]openai.MsgUsage{curr, usage})

	data, err := json.MarshalIndent(total, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal usage: %w", err)
	}
	if err := os.WriteFile(UsagePath, data, 0644); err != nil {
		return fmt.Errorf("save usage: %w", err)
	}
	return nil
}
