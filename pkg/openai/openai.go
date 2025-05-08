package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/responses"
)

func (m *Model) generate(
	r responses.ResponseNewParams, h func(string),
) ([]Msg, error) {
	response, err := m.cli.Responses.New(context.Background(), r)
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}
	if string(response.Error.Code) != "" {
		return nil, fmt.Errorf("openai: %s", response.Error.Message)
	}
	h(response.OutputText())
	return parse(response)
}

func (m *Model) stream(
	r responses.ResponseNewParams, h func(string),
) ([]Msg, error) {
	stream := m.cli.Responses.NewStreaming(context.Background(), r)
	defer stream.Close()

	// stream the response
	var response responses.Response
	for stream.Next() {
		data := stream.Current()
		if data.Response.Status == responses.ResponseStatusCompleted {
			response = data.AsResponseCompleted().Response
			break
		}
		h(parseStream(data))
	}
	println() // formatting

	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}
	if response.IncompleteDetails.Reason != "" {
		return nil, fmt.Errorf("openai: %s", response.IncompleteDetails.Reason)
	}
	return parse(&response)
}
