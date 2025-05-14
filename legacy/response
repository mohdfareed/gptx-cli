import "github.com/openai/openai-go/responses"

func parse(response *responses.Response) ([]MsgData, MsgUsage, error) {
	msgs := []MsgData{}
	for _, item := range response.Output {
		msg := MsgData{}
		switch item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			msgData := item.AsMessage().ToParam()
			msg = MsgData{OfOutputMessage: &msgData}
		case responses.ResponseFunctionWebSearch:
			msgData := item.AsWebSearchCall().ToParam()
			msg = MsgData{OfWebSearchCall: &msgData}
		case responses.ResponseFunctionToolCall:
			msgData := item.AsFunctionCall().ToParam()
			msg = MsgData{OfFunctionCall: &msgData}
		case responses.ResponseReasoningItem:
			msgData := item.AsReasoning().ToParam()
			msg = MsgData{OfReasoning: &msgData}
		}
		msgs = append(msgs, msg)
	}
	return msgs, response.Usage, nil
}

func (p *StreamParser) parseStream(event responses.ResponseStreamEventUnion) {
	switch variant := event.AsAny().(type) {
	case responses.ResponseTextDeltaEvent:
		p.Text <- variant.Delta
	case responses.ResponseRefusalDeltaEvent:
		p.Refusal <- variant.Delta
	case responses.ResponseFunctionCallArgumentsDeltaEvent:
		p.FunctionCall <- variant.Delta
	case responses.ResponseWebSearchCallSearchingEvent:
		p.WebSearch <- struct{}{}
	}
}

func (p *StreamParser) close() {
	close(p.Text)
	close(p.Refusal)
	close(p.FunctionCall)
	close(p.WebSearch)
}
