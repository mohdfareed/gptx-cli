package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/responses"
)

var WebSearchTool Tool = Tool{
	Name: "web_search",
	Desc: "search the web for information",
	Def: responses.ToolParamOfWebSearch(
		responses.WebSearchToolTypeWebSearchPreview,
	), // REVIEW: check preview status
	Init: func(config Config) error {
		return nil
	},
}

// Generate a reply from the model.
func (m *Model) Prompt(
	msgs []Msg, tools []Tool, h func(string),
) ([]Msg, error) {
	request := newRequest(m.config, msgs, tools)
	if m.config.Stream {
		return m.stream(request, h)
	} else {
		return m.generate(request, h)
	}
}

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
