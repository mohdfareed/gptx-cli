import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

// UserMsg creates a user message with the given text and files.
func UserMsg(text string, files []string) (MsgData, error) {
	var data []FileData
	for _, path := range files {
		file, err := readFile(path)
		if err != nil {
			return MsgData{}, fmt.Errorf("readFile: %w", err)
		}
		data = append(data, file)
	}

	if text != "" {
		data = append(data, FileData{
			OfInputText: &responses.ResponseInputTextParam{Text: text},
		})
	}

	msg := MsgData{
		OfInputMessage: &responses.ResponseInputItemMessageParam{
			Role: "user", Content: data,
		},
	}
	return msg, nil
}

func readFile(path string) (FileData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileData{}, fmt.Errorf("loadFile: %w", err)
	}

	switch filepath.Ext(path) {
	case ".jpg", ".jpeg", ".png", ".svg":
		return imageFile(data, path)
	default:
		return dataFile(data, path)
	}
}

func dataFile(data []byte, path string) (FileData, error) {
	// format := "# File: %s\n\n```%s\n%s\n```"
	// ext := filepath.Ext(path)
	// text := fmt.Sprintf(format, path, ext, string(data))
	// file := responses.ResponseInputTextParam{Text: text}
	// return FileData{OfInputText: &file}, nil

	file := responses.ResponseInputFileParam{
		FileID:   param.Opt[string]{Value: path},
		Filename: param.Opt[string]{Value: filepath.Base(path)},
		FileData: param.Opt[string]{Value: string(data)},
	} // REVIEW: convert to text message if not a text file
	return FileData{OfInputFile: &file}, nil
}

func imageFile(data []byte, path string) (FileData, error) {
	// infer MIME from extension; fallback to sniffing bytes
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	// base64‑encode and assemble
	b64 := base64.StdEncoding.EncodeToString(data)
	url := fmt.Sprintf("data:%s;base64,%s", mimeType, b64)

	// create image file
	image := responses.ResponseInputImageParam{
		FileID:   param.Opt[string]{Value: path},
		ImageURL: param.Opt[string]{Value: url},
	}
	return FileData{OfInputImage: &image}, nil
}
