// Package cfg handles configuration management.
package cfg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// AppName is the application name.
const AppName string = "gptx"

// AppName is the application name.
const EnvVarPrefix string = "GPTX_"

// AppDir is the user config directory for the app.
var AppDir string = func() string {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		return ""
	}
	return filepath.Join(configDir, AppName)
}()

// EnvMap returns the current environment variables as a map.
func EnvMap() map[string]string {
	m := make(map[string]string)
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 && parts[1] != "" {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

// LoadConfigFiles loads .gptx files in Git-like fashion:
// - Current directory and parent dirs (for project settings)
// - User's config directory (for global settings)
func LoadConfigFiles() {
	godotenv.Load(ConfigFiles()...)
}

// ConfigFiles returns paths to all relevant configuration files.
func ConfigFiles() []string {
	var files []string

	// Look for .gptx files in current directory and parent directories
	for dir, err := os.Getwd(); err == nil; dir = filepath.Dir(dir) {
		configFile := filepath.Join(dir, "."+AppName)
		if _, err := os.Stat(configFile); err == nil {
			files = append(files, configFile)
		}

		// Stop at root directory
		if dir == filepath.Dir(dir) {
			break
		}
	}

	// Add global config if it exists
	if AppDir != "" {
		globalConfig := filepath.Join(AppDir, "config")
		if _, err := os.Stat(globalConfig); err == nil {
			files = append(files, globalConfig)
		}
	}

	return files
}

// Auto-load configuration
func init() {
	LoadConfigFiles()
}
