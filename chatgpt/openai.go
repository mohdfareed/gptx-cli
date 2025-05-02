package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

// stream a response from the model.
func (m *Model) stream(message string) error {
	// generate stream request
	request := m.createRequest(message)
	stream := m.cli.Responses.NewStreaming(context.Background(), request)
	defer stream.Close()

	for { // stream the response
		for stream.Next() {
			data := stream.Current()
			println(parseStream(data))
		}

		if err := stream.Err(); err != nil {
			return fmt.Errorf("stream error: %w", err)
		}
	}
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
	prompt := param.Opt[string]{Value: message}
	params := responses.ResponseNewParams{
		Input:        responses.ResponseNewParamsInputUnion{OfString: prompt},
		Model:        m.config.Model,
		Instructions: param.Opt[string]{Value: m.config.SysPrompt},
		// Include: []responses.ResponseIncludable{
		// 	responses.ResponseIncludableComputerCallOutputOutputImageURL,
		// 	responses.ResponseIncludableFileSearchCallResults,
		// 	responses.ResponseIncludableMessageInputImageImageURL,
		// },
		// Reasoning: shared.ReasoningParam{
		// 	Effort:          shared.ReasoningEffortLow,
		// 	GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
		// },
	}
	// TODO: add attachments, tools, max tokens?, temp?, parallel calls?
	return params
}
