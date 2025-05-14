package gptx

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

// AppName is the name of the application.
const AppName string = "gptx"

// AppDir is the directory where application configuration files are stored.
var AppDir string = func() string {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		return ""
	}
	return filepath.Join(configDir, AppName)
}()

// EnvVar returns the environment variable name for a given field.
// It uses the struct field's "env" tag or the field name.
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

	// Format as GPTX_FIELD_NAME
	prefix := strings.ToUpper(AppName)
	postfix := strings.ToUpper(tag)
	return fmt.Sprintf("%s_%s", prefix, postfix)
}

// LoadConfigFiles loads configuration from .gptx files.
// Searches for files in:
// 1. Current directory and all parent directories
// 2. User's config directory
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
