package openai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/mohdfareed/gptx-cli/pkg/llm"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

func createMsgs(msgs []llm.Msg) []openAIMsg {
	openAImsgs := make([]openAIMsg, len(msgs))
	for i, msg := range msgs {

	}
}

func createOpenAIData(msg llm.Msg) openAIMsg {
	switch m := msg.(type) {
	case llm.TextMsg:
		return textMsg(m.Text)
	case llm.ToolResult:
		return toolMsg(m)
	}
}

func toolMsg(msg llm.ToolResult[any, any]) openAIMsg {
	output, err := json.Marshal(msg.Result)
	if err != nil {
		output = []byte(fmt.Sprintf("%v", msg.Result))
	}
	data := responses.ResponseInputItemFunctionCallOutputParam{
		CallID: msg.ID,
		Output: string(output),
	}
	return openAIMsg{OfFunctionCallOutput: &data}
}

func userMsg(msg llm.TextMsg) (openAIMsg, error) {
	content := []openAIData{}
	for _, file := range msg.Files {
		switch filepath.Ext(file.Path) {
		case ".png", ".jpg", ".jpeg", ".gif":
			data, err := imageFile(file.Content, file.Path)
			if err != nil {
				return openAIMsg{}, err
			}
			content = append(content, data)
		default:
			data, err := dataFile(file.Content, file.Path)
			if err != nil {
				return openAIMsg{}, err
			}
			content = append(content, data)
		}
	}

	data := responses.ResponseInputItemMessageParam{
		Role:    string(msg.Role),
		Content: content,
	}

	// TODO: handle
	// msg := openAIMsg{OfReasoning: &data}
	// msg := openAIMsg{OfInputMessage: &data}
	// msg := openAIMsg{OfOutputMessage: &data}
	// msg := openAIMsg{OfWebSearchCall: &data}
	// msg := openAIMsg{OfFunctionCall: &data}
	// msg := openAIMsg{OfFunctionCallOutput: &data}
	return data, nil
}

// MARK: Files
// ============================================================================

func textMsg(data string) openAIData {
	file := responses.ResponseInputTextParam{Text: data}
	return openAIData{OfInputText: &file}
}

// func textFile(data []byte, path string) openAIFile {
// 	format := "# File: %s\n\n```%s\n%s\n```"
// 	ext := filepath.Ext(path)
// 	text := fmt.Sprintf(format, path, ext, string(data))
// 	return textMsg(text)
// }

func dataFile(data []byte, path string) (openAIData, error) {
	file := responses.ResponseInputFileParam{
		FileID:   param.Opt[string]{Value: path},
		Filename: param.Opt[string]{Value: filepath.Base(path)},
		FileData: param.Opt[string]{Value: string(data)},
	} // TODO: test, convert to text message with file block
	return openAIData{OfInputFile: &file}, nil
}

func imageFile(data []byte, path string) (openAIData, error) {
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
	return openAIData{OfInputImage: &image}, nil
}
