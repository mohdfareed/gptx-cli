package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

func Generate(
	c openai.Client, r ModelRequest, h func(string),
) ([]MsgData, MsgUsage, error) {
	response, err := c.Responses.New(context.Background(), r)
	if err != nil {
		return nil, response.Usage, fmt.Errorf("openai: %w", err)
	}
	if string(response.Error.Code) != "" {
		msg := response.Error.Message
		return nil, response.Usage, fmt.Errorf("openai: %s", msg)
	}
	h(response.OutputText())
	return parse(response)
}

func Stream(
	c openai.Client, r ModelRequest, h func(string),
) ([]MsgData, MsgUsage, error) {
	stream := c.Responses.NewStreaming(context.Background(), r)
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
		return nil, response.Usage, fmt.Errorf("openai: %w", err)
	}
	if response.IncompleteDetails.Reason != "" {
		reason := response.IncompleteDetails.Reason
		return nil, response.Usage, fmt.Errorf("openai: %s", reason)
	}
	return parse(&response)
}
