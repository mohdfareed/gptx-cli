package files

import (
	"fmt"
	"os"
	"path/filepath"
)

// AttachedFile represents a file to be sent with a model message.
type AttachedFile struct {
	Name    string // Base name of the file (no directory paths)
	Path    string // Full filesystem path to the file
	Content string // Entire contents of the file
}

// LoadFiles scans the given glob patterns and loads all matching files.
// It ignores directories and returns a slice of AttachedFiles or an error.
func LoadFiles(patterns []string) ([]AttachedFile, error) {
	var result []AttachedFile

	for _, pat := range patterns {
		matches, err := filepath.Glob(pat)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", pat, err)
		}

		for _, path := range matches {
			info, err := os.Stat(path)
			if err != nil {
				continue // skip unreadable paths
			}
			if info.IsDir() {
				continue // skip directories
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("reading file %q: %w", path, err)
			}

			result = append(result, AttachedFile{
				Name:    filepath.Base(path),
				Path:    path,
				Content: string(data),
			})
		}
	}

	return result, nil
}
