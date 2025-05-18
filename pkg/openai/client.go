package openai

import (
	"context"
	"fmt"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

type OpenAIClient struct {
	client openai.Client
	userID string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

func (c *OpenAIClient) Generate(
	ctx context.Context, model gptx.Model, prompt string,
) error {
	msg, err := UserMsg(prompt, model.Config.Files)
	if err != nil {
		return fmt.Errorf("openai: %w", err)
	}
	msgs := []MsgData{msg}

	req := NewRequest(model, msgs, shared.ReasoningEffortHigh, c.userID)
	stream := c.client.Responses.NewStreaming(ctx, req)
	model.Events.Start.Emit(ctx, model.Config)
	defer stream.Close()

	// stream the response
	var response responses.Response
	for stream.Next() {
		data := stream.Current()
		if data.Response.Status == responses.ResponseStatusCompleted {
			response = data.AsResponseCompleted().Response
			break
		}
		parseStream(data, *model.Events, ctx)
	}

	// check for errors
	if err := stream.Err(); err != nil {
		model.Events.Done.Emit(ctx, response.Usage.RawJSON())
		return fmt.Errorf("openai: %w", err)
	}
	if response.IncompleteDetails.Reason != "" {
		reason := response.IncompleteDetails.Reason
		model.Events.Done.Emit(ctx, response.Usage.RawJSON())
		return fmt.Errorf("openai: %s", reason)
	}

	// parse the response
	_, usage, err := parse(&response, *model.Events, ctx)
	model.Events.Done.Emit(ctx, usage.RawJSON())
	return err
}
