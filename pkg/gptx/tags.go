package gptx

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const tagRegex = `@(.+?)\((.*?)\)` // @tag(args)
var regex = regexp.MustCompile(tagRegex)

// ProcessTags processes special tags in the prompt text and replaces them
// with their content. Currently supports:
// - @file(path) - Includes the content of the entire file
// - @file(path:start-end) - Includes specific line range from a file
func ProcessTags(prompt string) (string, error) {
	// parse tags
	matches := regex.FindAllStringSubmatch(prompt, -1)
	if len(matches) == 0 {
		return prompt, nil
	}

	// process each tag
	for _, match := range matches {
		var result string
		var err error

		switch match[1] {
		case "file":
			result, err = fileTag(match[2])
		default: // unknown tag
			continue
		}

		if err != nil {
			return "", fmt.Errorf("tags: %w", err)
		}
		prompt = strings.Replace(prompt, match[0], result, 1)
	}
	return prompt, nil
}

// MARK: File Tags
// ============================================================================

const fileTagRegex = `^(.*?)(?::(\d+)-(\d+))?$` // ^filepath[:start-end]$
var fileRegex = regexp.MustCompile(fileTagRegex)

func fileTag(args string) (string, error) {
	// parse the file tag
	match := fileRegex.FindStringSubmatch(args)
	if len(match) == 0 {
		return "", fmt.Errorf("invalid file: %q", args)
	}
	path := match[1]

	// read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	file := strings.Split(string(data), "\n")
	id := path

	// parse the start and end lines
	if len(match[2]) > 0 && len(match[3]) > 0 {
		start, err := strconv.Atoi(match[2])
		if err != nil {
			return "", fmt.Errorf("file tag: %w", err)
		}
		end, err := strconv.Atoi(match[3])
		if err != nil {
			return "", fmt.Errorf("file tag: %w", err)
		}
		if start < 0 || end < 0 || start > end {
			return "", fmt.Errorf("file tag: %d-%d", start, end)
		}
		if start > len(file) {
			return "", fmt.Errorf("file tag: start line %d out of range", start)
		}
		endLine := end
		if endLine > len(file) {
			endLine = len(file)
		}
		file = file[start-1 : endLine]
		id = fmt.Sprintf("%s:%d-%d", id, start, end)
	}

	// create tag block
	tag := "\nFile: %s\n\n```%s\n%s\n```\n"
	ext := filepath.Ext(path)
	text := strings.Join(file, "\n")
	return fmt.Sprintf(tag, id, ext, text), nil
}
