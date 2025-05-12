package gptx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// AppName is the name of the application.
const AppName string = "gptx"

// AppDir is the directory where the application configuration files are stored.
var AppDir string = func() string {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		return ""
	}
	return filepath.Join(configDir, AppName)
}()

// EnvVar returns the environment variable name for a given key.
func EnvVar(name string) string {
	prefix := strings.ToUpper(AppName)
	return fmt.Sprintf("%s_%s", prefix, name)
}

// LoadConfigFiles loads configuration from dotenv files.
// It searches hierarchically from the current directory up to the root,
// following Git-like behavior for .gptx files.
func LoadConfigFiles() {
	godotenv.Load(ConfigFiles()...)
}

// ConfigFiles returns the paths of configuration files to load.
// It searches for:
// - .gptx files in the current directory and all parent directories
// - config file in the application directory
func ConfigFiles() []string {
	var files []string

	// Look for .gptx files in current directory and all parents
	for dir, err := os.Getwd(); err == nil; dir = filepath.Dir(dir) {
		f := filepath.Join(dir, "."+AppName)
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}

		// Stop at root directory
		if dir == filepath.Dir(dir) {
			break
		}
	}

	// Global application config
	if AppDir != "" {
		globalConfig := filepath.Join(AppDir, "config")
		if _, err := os.Stat(globalConfig); err == nil {
			files = append(files, globalConfig)
		}
	}

	return files
}

func init() {
	LoadConfigFiles()
}
