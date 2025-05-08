package gptx

// import (
// 	"context"
// 	"fmt"

// 	"github.com/openai/openai-go"
// 	"github.com/openai/openai-go/option"
// 	"github.com/urfave/cli/v3"
// )

// // The OpenAI client.
// type Model struct {
// 	cli    openai.Client
// 	config Config
// 	tools  ModelTools
// }

// // Create a new OpenAI model client.
// func CreateModel(config Config, tools ModelTools) Model {
// 	client := openai.NewClient(option.WithAPIKey(config.APIKey))
// 	return Model{
// 		cli:    client,
// 		config: config,
// 		tools:  tools,
// 	}
// }

// // Validate the model configuration by generating a request with no messages.
// // It overrides the max tokens to 0 to avoid generating a response.
// func (m *Model) Validate() error {
// 	tmp := m.config.MaxTokens
// 	m.config.MaxTokens = 16
// 	_, err := m.Prompt([]Msg{UserMsg("")}, []ModelTool{}, func(string) {})
// 	m.config.MaxTokens = tmp
// 	if err != nil {
// 		return fmt.Errorf("validation: %w", err)
// 	}
// 	return nil
// }

// // MARK: Model Pipeline
// // ============================================================================

// // Generate or stream a model reply.
// func (m *Model) Run(message string, handler func(string)) (string, error) {
// 	// process message tags
// 	message, err := ProcessTags(message)
// 	if err != nil {
// 		return "", fmt.Errorf("model: %w", err)
// 	}

// 	// load attachments
// 	files, err := MsgFiles(m.config.Files)
// 	if err != nil {
// 		return "", fmt.Errorf("model: %w", err)
// 	}

// 	// initialize tools
// 	err = m.tools.Init() // initialize the model tools
// 	if err != nil {
// 		return "", fmt.Errorf("model: %w", err)
// 	}
// 	tools := m.tools.enabledTools()

// 	// load the chat history
// 	history, err := LoadChat(m.config.Chat)
// 	if err != nil {
// 		return "", fmt.Errorf("history: %w", err)
// 	}

// 	// create prompt
// 	msgs := []Msg{}
// 	if len(files) > 0 {
// 		msgs = append(msgs, FilesMsg(files))
// 	}
// 	if message != "" {
// 		history.Add(UserMsg(message))
// 	} // ignore if re-prompting without a new message
// 	msgs = append(msgs, history.Msgs...)

// 	// prompt the model
// 	reply, err := m.Prompt(msgs, tools, handler)
// 	if err != nil {
// 		return "", fmt.Errorf("model: %w", err)
// 	}

// 	// save to the chat history
// 	for _, msg := range reply {
// 		history.Add(msg)
// 		if err := history.Save(); err != nil {
// 			return "", fmt.Errorf("history: %w", err)
// 		}
// 	}

// 	// TODO: use tool if reply is function call and rerun with result

// 	// return the model reply
// 	usage := fmt.Sprintf(
// 		"\n\n%s",
// 		reply[len(reply)-1].Usage.RawJSON(),
// 	)
// 	return usage, nil
// }

// // MARK: CLI
// // ============================================================================

// func ValidateCMD(model *Model) *cli.Command {
// 	return &cli.Command{
// 		Name: "check", Usage: "validate the model config",
// 		Action: func(ctx context.Context, cmd *cli.Command) error {
// 			ToolsCMD(model.tools)
// 			println(Theme.Magenta + "validating model config..." + Theme.Reset)
// 			if err := model.Validate(); err != nil {
// 				println(Theme.Error + err.Error())
// 			} else {
// 				println(Theme.Success + "config valid" + Theme.Reset)
// 			}
// 			return nil
// 		},
// 	}
// }

// func MsgCMD(model *Model) *cli.Command {
// 	var msgs []string
// 	return &cli.Command{
// 		Name:  AppName,
// 		Usage: "message an OpenAI model",
// 		Arguments: []cli.Argument{
// 			&cli.StringArgs{
// 				Name:        "msgs",
// 				UsageText:   "the message to send to the model",
// 				Destination: &msgs,
// 				Max:         -1,
// 			},
// 		},
// 		Action: func(ctx context.Context, cmd *cli.Command) error {
// 			prompt, err := PromptUser(model.config, msgs)
// 			if err != nil {
// 				println(Theme.Error + err.Error())
// 			}

// 			isStreamed := model.config.Stream
// 			reply, err := model.Run(prompt, func(c string) {
// 				print(c)
// 			})

// 			if !isStreamed {
// 				print("@model: " + reply)
// 			}
// 			return err
// 		},
// 	}
// }
