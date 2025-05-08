package openai

import (
	"github.com/openai/openai-go/responses"
)

// MsgUsage represents the usage of a model response.
type MsgUsage = responses.ResponseUsage

// TotalUsage calculates the total usage from a slice of MsgUsage.
func TotalUsage(usages []MsgUsage) MsgUsage {
	total := MsgUsage{}
	for _, usage := range usages {
		total.InputTokens += usage.InputTokens
		total.OutputTokens += usage.OutputTokens
		total.TotalTokens += usage.TotalTokens

		reasoning := usage.OutputTokensDetails.ReasoningTokens
		total.OutputTokensDetails.ReasoningTokens += reasoning
		cache := usage.InputTokensDetails.CachedTokens
		total.InputTokensDetails.CachedTokens += cache
	}
	return total
}
