package gptx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const AppName string = "gptx"

var AppDir string = func() string {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		return ""
	}
	return filepath.Join(configDir, AppName)
}()

func EnvVar(name string) string {
	prefix := strings.ToUpper(AppName)
	return fmt.Sprintf("%s_%s", prefix, name)
}
