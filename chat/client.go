package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

// The OpenAI client.
type Client struct {
	cli   *openai.Client
	model string
	instr string
}

// Create a new OpenAI Client.
func NewClient(apiKey, model string) *Client {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &Client{
		cli:   &client,
		model: model,
	}
}

// Prompt the model with a new message.
func (c *Client) Prompt(message string) (string, error) {
	// construct request
	prompt := param.Opt[string]{Value: message}
	instructs := param.Opt[string]{Value: c.instr}
	params := responses.ResponseNewParams{
		Input:        responses.ResponseNewParamsInputUnion{OfString: prompt},
		Model:        c.model,
		Instructions: instructs,
	}

	// generate response
	response, err := c.cli.Responses.New(context.Background(), params)
	return response.OutputText(), err
}
