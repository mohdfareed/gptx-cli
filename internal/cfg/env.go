// Package cfg handles configuration management.
package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

// AppName is the application name.
const AppName string = "gptx"

// AppDir is the user config directory for the app.
var AppDir string = func() string {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		return ""
	}
	return filepath.Join(configDir, AppName)
}()

// EnvVar returns the env var name for a config field.
// Example: For field "APIKey" with tag `env:"api_key"`, returns "GPTX_API_KEY"
func EnvVar(obj *Config, field string) string {
	var tag string
	if obj != nil {
		// Use reflection to get the field tag
		t := reflect.TypeOf(*obj)
		f, found := t.FieldByName(field)
		if !found {
			panic(fmt.Sprintf("field '%s' not found in type Config", field))
		}
		tag = f.Tag.Get("env")
	} else {
		tag = field // Use the field name directly
	}

	// Format as GPTX_FIELD_NAME (e.g., GPTX_API_KEY)
	prefix := strings.ToUpper(AppName)
	postfix := strings.ToUpper(tag)
	return fmt.Sprintf("%s_%s", prefix, postfix)
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
