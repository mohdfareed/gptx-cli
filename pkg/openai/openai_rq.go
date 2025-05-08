package openai

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// MARK: Requests
// ============================================================================

func newRequest(
	config Config, msgs []Msg, tools []Tool,
) responses.ResponseNewParams {
	// config
	data := responses.ResponseNewParams{
		Model:        config.Model,
		Instructions: param.Opt[string]{Value: config.SysPrompt},
		Include: []responses.ResponseIncludable{
			responses.ResponseIncludableMessageInputImageImageURL,
		},
		MaxOutputTokens:   param.Opt[int64]{Value: config.MaxTokens},
		ParallelToolCalls: param.Opt[bool]{Value: true},
		Temperature:       param.Opt[float64]{Value: config.Temp},
		User:              param.Opt[string]{Value: userID},
	}

	// reasoning
	if config.ReasonEffort != "" {
		data.Reasoning = shared.ReasoningParam{
			Effort:          config.ReasonEffort,
			GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
		}
	}

	// tools
	var toolDefs []ToolDef
	for _, tool := range tools {
		toolDefs = append(toolDefs, tool.Def)
	}
	data.Tools = toolDefs

	// messages
	reqMsgs := make([]MsgData, len(msgs))
	for i, msg := range msgs {
		reqMsgs[i] = msg.Data
	}
	data.Input.OfInputItemList = reqMsgs
	return data
}

// MARK: Response
// ============================================================================

func parseStream(event responses.ResponseStreamEventUnion) string {
	return event.Delta
}

func parse(response *responses.Response) ([]Msg, error) {
	msgs := []Msg{}
	for _, item := range response.Output {
		msg := Msg{}
		switch item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			msgData := item.AsMessage().ToParam()
			msg.Data = MsgData{OfOutputMessage: &msgData}
		case responses.ResponseFunctionWebSearch:
			msgData := item.AsWebSearchCall().ToParam()
			msg.Data = MsgData{OfWebSearchCall: &msgData}
		case responses.ResponseFunctionToolCall:
			msgData := item.AsFunctionCall().ToParam()
			msg.Data = MsgData{OfFunctionCall: &msgData}
		case responses.ResponseReasoningItem:
			msgData := item.AsReasoning().ToParam()
			msg.Data = MsgData{OfReasoning: &msgData}
		case responses.ResponseComputerToolCall:
			msgData := item.AsComputerCall().ToParam()
			msg.Data = MsgData{OfComputerCall: &msgData}
		case responses.ResponseFileSearchToolCall:
			msgData := item.AsFileSearchCall().ToParam()
			msg.Data = MsgData{OfFileSearchCall: &msgData}
		}
		msgs = append(msgs, msg)
	}

	// set usage of last message
	if len(msgs) > 0 {
		msgs[len(msgs)-1].Usage = &response.Usage
	}
	return msgs, nil
}

// MARK: Messages
// ============================================================================

type File = responses.ResponseInputContentUnionParam

// TextMsg creates a user message with the given text.
func UserMsg(text string) Msg {
	msg := responses.ResponseInputTextParam{Text: text}
	data := MsgData{
		OfInputMessage: &responses.ResponseInputItemMessageParam{
			Role: "user", Content: []File{{OfInputText: &msg}},
		},
	}
	return Msg{Data: data}
}

// FilesMsg creates a user message with the given files.
func FilesMsg(files []File) Msg {
	data := MsgData{
		OfInputMessage: &responses.ResponseInputItemMessageParam{
			Role: "user", Content: files,
		},
	}
	return Msg{Data: data}
}

func TextFile(data []byte, path string) (File, error) {
	format := "# File: %s\n\n```%s\n%s\n```"
	ext := filepath.Ext(path)
	text := fmt.Sprintf(format, path, ext, string(data))
	file := responses.ResponseInputTextParam{Text: text}
	return File{OfInputText: &file}, nil
}

func DataFile(data []byte, path string) (File, error) {
	file := responses.ResponseInputFileParam{
		FileID:   param.Opt[string]{Value: path},
		Filename: param.Opt[string]{Value: filepath.Base(path)},
		FileData: param.Opt[string]{Value: string(data)},
	} // TODO: test, convert to text message with file block
	return File{OfInputFile: &file}, nil
}

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
