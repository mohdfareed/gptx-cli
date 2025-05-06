package openai

import (
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Create a new OpenAI model client.
func CreateModel(
	config ModelConfig, tools ModelTools, platform *ModelPlatform[any],
) Model {
	client := openai.NewClient(option.WithAPIKey(config.APIKey))
	return Model{
		cli:    client,
		config: config,
		tools:  tools,
	}
}

// Validate the model configuration by generating a request with no messages.
// It overrides the max tokens to 0 to avoid generating a response.
func (m *Model) Validate() error {
	tmp := m.config.MaxTokens
	m.config.MaxTokens = 16
	_, err := m.Prompt([]openAIMsg{UserMsg("")}, []Tool{}, func(string) {})
	m.config.MaxTokens = tmp
	if err != nil {
		return fmt.Errorf("validation: %w", err)
	}
	return nil
}

// MARK: Model Pipeline
// ============================================================================

// Generate or stream a model reply.
func (m *Model) Msg(chat openAIChat, handler func(string)) (string, error) {
	// load attachments
	files, err := MsgFiles(m.config.Files)
	if err != nil {
		return "", fmt.Errorf("model: %w", err)
	}

	// initialize tools
	err = m.tools.Init() // initialize the model tools
	if err != nil {
		return "", fmt.Errorf("model: %w", err)
	}
	tools := m.tools.enabledTools()

	// load the chat history
	history, err := LoadChat(m.config.Chat)
	if err != nil {
		return "", fmt.Errorf("history: %w", err)
	}

	// create prompt
	msgs := []openAIMsg{}
	if len(files) > 0 {
		msgs = append(msgs, FilesMsg(files))
	}
	if message != "" {
		history.Add(UserMsg(message))
	} // ignore if re-prompting without a new message
	msgs = append(msgs, history.Msgs...)

	// prompt the model
	reply, err := m.Prompt(msgs, tools, handler)
	if err != nil {
		return "", fmt.Errorf("model: %w", err)
	}

	// save to the chat history
	for _, msg := range reply {
		history.Add(msg)
		if err := history.Save(); err != nil {
			return "", fmt.Errorf("history: %w", err)
		}
	}

	// TODO: use tool if reply is function call and rerun with result

	// return the model reply
	usage := fmt.Sprintf(
		"\n\n%s",
		reply[len(reply)-1].Usage.RawJSON(),
	)
	return usage, nil
}
