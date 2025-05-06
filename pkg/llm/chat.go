package llm

// MARK: Messages
// ============================================================================

type MsgSender string

const (
	MsgSenderUser   MsgSender = "user"
	MsgSenderModel  MsgSender = "model"
	MsgSenderSystem MsgSender = "system"
	MsgSenderDev    MsgSender = "developer"
)

// Msg represents a single turn in the conversation.
type Msg interface {
	Sender() MsgSender
	UsageData() *Usage
}

// Chat holds a chat context.
type Chat struct {
	Msgs []Msg `json:"messages"`
}

// MARK: Messages
// ============================================================================

// File is a file attachment to a message.
type File struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
}

// TextMsg represents a text message with attachments.
type TextMsg struct {
	Role  MsgSender `json:"role"`
	Text  string    `json:"text"`
	Files []File    `json:"files,omitempty"`
	Usage Usage     `json:"usage,omitempty"`
	Data  any       `json:"data,omitempty"`
}

// ToolCall is a model tool's call.
type ToolCall[Params any] struct {
	ID    string `json:"id"`
	Tool  Tool   `json:"tool"`
	Args  Params `json:"args"`
	Usage Usage  `json:"usage,omitempty"`
}

// ToolResult is a model tool's call results.
type ToolResult[Params any, Results any] struct {
	ID     string  `json:"args"`
	Result Results `json:"results"`
}

func (m TextMsg) Sender() MsgSender                     { return m.Role }
func (m TextMsg) UsageData() *Usage                     { return &m.Usage }
func (m ToolCall[Params]) Sender() MsgSender            { return MsgSenderDev }
func (m ToolCall[Params]) UsageData() *Usage            { return &m.Usage }
func (m ToolResult[Params, Results]) Sender() MsgSender { return MsgSenderDev }
func (m ToolResult[Params, Results]) UsageData() *Usage { return nil }

// MARK: Chat
// ============================================================================

func CreateChat(msgs []Msg, model *Model) Chat {
	chat := Chat{Msgs: msgs}
	model.Finished.Subscribe(func(msg Msg) {
		chat.Add(msg)
	})
	model.ToolResult.Subscribe(func(msg ToolResult[any, any]) {
		chat.Add(msg)
	})
	return chat
}

// Add appends a message to the chat.
func (c *Chat) Add(msg Msg) {
	c.Msgs = append(c.Msgs, msg)
}
