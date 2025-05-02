package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// stream a response from the model.
func (m *Model) stream(message string) error {
	// generate stream request
	request := m.createRequest(message)
	stream := m.cli.Responses.NewStreaming(context.Background(), request)
	defer stream.Close()

	// stream the response
	for stream.Next() {
		data := stream.Current()
		print(parseStream(data))
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	println()
	return nil
}

// generate a response from the model.
func (m *Model) generate(message string) error {
	request := m.createRequest(message)
	response, err := m.cli.Responses.New(context.Background(), request)
	println(parse(response))
	return err
}

// MARK: Response
// ============================================================================

func parseStream(response responses.ResponseStreamEventUnion) string {
	// TODO: add support for other response data
	return response.Delta
}

func parse(response *responses.Response) string {
	// TODO: add support for other response data
	return response.OutputText()
}

// MARK: Request
// ============================================================================

func (m *Model) createRequest(message string) responses.ResponseNewParams {
	params := responses.ResponseNewParams{
		Model:        m.config.Model,
		Instructions: param.Opt[string]{Value: m.config.SysPrompt},
		Include: []responses.ResponseIncludable{
			responses.ResponseIncludableMessageInputImageImageURL,
		},
	}

	params.Tools = []responses.ToolUnionParam{
		responses.ToolParamOfWebSearch(
			responses.WebSearchToolTypeWebSearchPreview,
		),
	}

	if m.config.Model[0] == 'o' {
		params.Reasoning = shared.ReasoningParam{
			Effort:          shared.ReasoningEffortMedium,
			GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
		}
	}

	params.Input.OfInputItemList = []responses.ResponseInputItemUnionParam{
		{
			OfInputMessage: &responses.ResponseInputItemMessageParam{
				Role: string(responses.ResponseInputMessageItemRoleUser),
				Content: responses.ResponseInputMessageContentListParam{
					{
						OfInputText: &responses.ResponseInputTextParam{
							Text: message,
						},
					},
				},
			},
		},
	}

	// TODO: add attachments, custom tools (shell)
	// TODO: configs: max tokens?, temp?, parallel calls?
	return params
}
