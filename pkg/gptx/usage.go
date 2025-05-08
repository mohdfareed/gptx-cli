package gptx

// import (
// 	"context"
// 	"fmt"
// 	"strconv"

// 	"github.com/openai/openai-go/responses"
// 	"github.com/urfave/cli/v3"
// )

// type MsgUsage = responses.ResponseUsage

// // Return the total aggregated usage of the history.
// func TotalUsage(msgs []Msg) MsgUsage {
// 	usage := MsgUsage{}
// 	for _, msg := range msgs {
// 		usage.InputTokens += msg.Usage.InputTokens
// 		usage.OutputTokens += msg.Usage.OutputTokens
// 		usage.TotalTokens += msg.Usage.TotalTokens

// 		reasoning := msg.Usage.OutputTokensDetails.ReasoningTokens
// 		usage.OutputTokensDetails.ReasoningTokens += reasoning
// 		cache := msg.Usage.InputTokensDetails.CachedTokens
// 		usage.InputTokensDetails.CachedTokens += cache
// 	}
// 	return usage
// }

// // MARK: CLI
// // ============================================================================

// func UsageCMD(config Config) *cli.Command {
// 	printRow := func(key string, value string) {
// 		println(
// 			Theme.Dim + key + Theme.Reset + " " +
// 				Theme.Bold + value + Theme.Reset,
// 		)
// 	}

// 	return &cli.Command{
// 		Name: "usage", Usage: "show the chat's usage",
// 		Action: func(ctx context.Context, cmd *cli.Command) error {
// 			history, err := LoadChat(config.Chat)
// 			if err != nil {
// 				return fmt.Errorf("usage: %w", err)
// 			}
// 			last := history.Last(1)[0].Usage
// 			total := TotalUsage(history.Msgs)

// 			println(Theme.Bold + "last tokens usage:" + Theme.Reset)
// 			lastReasoning := last.OutputTokensDetails.ReasoningTokens
// 			lastCache := last.InputTokensDetails.CachedTokens
// 			printRow(Theme.Red+"    input:"+Theme.Reset, strconv.Itoa(int(last.InputTokens)))
// 			printRow(Theme.Green+"   output:"+Theme.Reset, strconv.Itoa(int(last.OutputTokens)))
// 			printRow(Theme.Blue+"    total:"+Theme.Reset, strconv.Itoa(int(last.TotalTokens)))
// 			printRow("reasoning:", strconv.Itoa(int(lastReasoning)))
// 			printRow("   cached:", strconv.Itoa(int(lastCache)))

// 			println(Theme.Bold + "total tokens usage:" + Theme.Reset)
// 			totalReasoning := total.OutputTokensDetails.ReasoningTokens
// 			totalCache := total.InputTokensDetails.CachedTokens
// 			printRow(Theme.Red+"    input:"+Theme.Reset, strconv.Itoa(int(total.InputTokens)))
// 			printRow(Theme.Green+"   output:"+Theme.Reset, strconv.Itoa(int(total.OutputTokens)))
// 			printRow(Theme.Blue+"    total:"+Theme.Reset, strconv.Itoa(int(total.TotalTokens)))
// 			printRow("reasoning:", strconv.Itoa(int(totalReasoning)))
// 			printRow("   cached:", strconv.Itoa(int(totalCache)))
// 			return nil
// 		},
// 	}
// }
