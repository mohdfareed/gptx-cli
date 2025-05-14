package files

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const tagRegex = `^@(.+)\((.*)\)$` // ^@tag(args)$
var regex = regexp.MustCompile(tagRegex)

func ProcessTags(prompt string) (string, []string, error) {
	// parse tags
	matches := regex.FindAllStringSubmatch(prompt, -1)
	if len(matches) == 0 {
		return prompt, nil, nil
	}

	// process each tag
	var attachments []string
	for _, match := range matches {
		var result string
		switch match[1] {
		case "file":
			r, a, err := fileTag(match[2])
			if err != nil {
				return "", nil, fmt.Errorf("file tag: %w", err)
			}
			attachments = append(attachments, a...)
			result = r
		default: // unknown tag
			continue
		}
		prompt = strings.Replace(prompt, match[0], result, 1)
	}
	return prompt, attachments, nil
}

// MARK: File Tags
// ============================================================================

const fileTagRegex = `^(.*?)(?::(\d+)-(\d+))?$` // ^filepath[:start-end]$
var fileRegex *regexp.Regexp = regexp.MustCompile(fileTagRegex)

func fileTag(args string) (string, []string, error) {
	// parse the file tag
	match := fileRegex.FindStringSubmatch(args)
	if len(match) == 0 {
		return "", nil, fmt.Errorf("invalid file: %q", args)
	}
	path := match[1]

	// read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	file := strings.Split(string(data), "\n")
	id := path

	// parse the start and end lines
	if len(match) == 3 {
		start, err := strconv.Atoi(match[2])
		if err != nil {
			return "", nil, fmt.Errorf("start line: %w", err)
		}
		end, err := strconv.Atoi(match[3])
		if err != nil {
			return "", nil, fmt.Errorf("end line: %w", err)
		}
		if start < 0 || end < 0 || start > end {
			return "", nil, fmt.Errorf("invalid range: %d-%d", start, end)
		}
		file = file[start-1 : end-1]
		id = fmt.Sprintf("%s:%d-%d", id, start, end)
	}

	// create tag block
	tag := "\nFile: %s\n\n```%s\n%s\n```\n"
	ext := filepath.Ext(path)
	text := strings.Join(file, "\n")
	return fmt.Sprintf(tag, id, ext, text), nil, nil
}
