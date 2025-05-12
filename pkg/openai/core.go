package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

// Generate sends a request to the OpenAI API and returns the generated
// messages and usage information. It uses the provided [StreamParser] to
// handle streaming events.
func Generate(
	c openai.Client, r ModelRequest, p StreamParser,
) ([]MsgData, MsgUsage, error) {
	stream := c.Responses.NewStreaming(context.Background(), r)
	defer stream.Close()
	defer p.close()

	// stream the response
	var response responses.Response
	for stream.Next() {
		data := stream.Current()
		if data.Response.Status == responses.ResponseStatusCompleted {
			response = data.AsResponseCompleted().Response
			break
		}
		p.parseStream(data)
	}

	// check for errors
	if err := stream.Err(); err != nil {
		return nil, response.Usage, fmt.Errorf("openai: %w", err)
	}
	if response.IncompleteDetails.Reason != "" {
		reason := response.IncompleteDetails.Reason
		return nil, response.Usage, fmt.Errorf("openai: %s", reason)
	}
	return parse(&response)
}
