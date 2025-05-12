package main

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/mohdfareed/gptx-cli/pkg/gptx"
)

func init() {
	godotenv.Load(configFIles()...)
}

func configFIles() []string {
	var files []string // cwd, parents, app dir

	// load files from cwd then its parents
	for dir, _ := os.Getwd(); ; dir = filepath.Dir(dir) {
		f := filepath.Join(dir, "."+gptx.AppName)
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}

		if dir == filepath.Dir(dir) {
			break // reached root
		}
	}

	// support $XDG_CONFIG_HOME and %APPDATA%
	if gptx.AppDir != "" {
		f := filepath.Join(gptx.AppDir, "config")
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}
	}
	return files
}
