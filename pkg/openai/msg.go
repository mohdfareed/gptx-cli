package openai

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

// MsgData represents the data structure for a message.
type MsgData = responses.ResponseInputItemUnionParam

// File represents the data structure for a file attachment.
type File = responses.ResponseInputContentUnionParam

// UserMsg creates a user message with the given text and files.
func UserMsg(text string, files []File) MsgData {
	if text != "" {
		files = append(files, File{
			OfInputText: &responses.ResponseInputTextParam{Text: text},
		})
	}

	data := MsgData{
		OfInputMessage: &responses.ResponseInputItemMessageParam{
			Role: "user", Content: files,
		},
	}
	return data
}

// TextFile creates a text file with the given data and path.
func TextFile(data []byte, path string) (File, error) {
	format := "# File: %s\n\n```%s\n%s\n```"
	ext := filepath.Ext(path)
	text := fmt.Sprintf(format, path, ext, string(data))
	file := responses.ResponseInputTextParam{Text: text}
	return File{OfInputText: &file}, nil
}

// DataFile creates a data file with the given data and path.
func DataFile(data []byte, path string) (File, error) {
	file := responses.ResponseInputFileParam{
		FileID:   param.Opt[string]{Value: path},
		Filename: param.Opt[string]{Value: filepath.Base(path)},
		FileData: param.Opt[string]{Value: string(data)},
	} // TODO: test, convert to text message with file block
	return File{OfInputFile: &file}, nil
}

// ImageFile creates an image file with the given data and path.
func ImageFile(data []byte, path string) (File, error) {
	// infer MIME from extension; fallback to sniffing bytes
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	// base64‑encode and assemble
	b64 := base64.StdEncoding.EncodeToString(data)
	url := fmt.Sprintf("data:%s;base64,%s", mimeType, b64)

	image := responses.ResponseInputImageParam{
		FileID:   param.Opt[string]{Value: path},
		ImageURL: param.Opt[string]{Value: url},
	}
	return File{OfInputImage: &image}, nil
}
