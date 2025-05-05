package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// MsgFiles gets message files from the given paths.
func MsgFiles(paths []string) ([]File, error) {
	_resolved, err := resolveFiles(paths)
	if err != nil {
		return nil, fmt.Errorf("files: %w", err)
	}
	paths = _resolved

	files := make([]File, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("file %q: %w", path, err)
		}

		var fileData File
		switch filepath.Ext(path) {
		// image files
		case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".tiff", ".svg":
			image, err := ImageFile(data, path)
			if err != nil {
				return nil, fmt.Errorf("image file %q: %w", path, err)
			}
			fileData = image

		// text files
		default: // TODO: test and replace with textFile
			data, err := DataFile(data, path)
			if err != nil {
				return nil, fmt.Errorf("data file %q: %w", path, err)
			}
			fileData = data
		}
		files = append(files, fileData)
	}
	return files, nil
}

// MARK: Helpers
// ============================================================================

func resolveFiles(paths []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("pattern %q: %w", path, err)
		}

		if len(matches) == 0 {
			warn(fmt.Errorf("not found: %q", path))
			return []string{}, nil
		}
		files = append(files, matches...)
	}
	return files, nil
}
