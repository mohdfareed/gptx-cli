package main

import (
	"fmt"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
)

// MARK: Views
// ============================================================================

// func printMsg(msg Msg) {
// 	// msg.
// }

// MARK: Model Prefix
// ============================================================================

func modelPrefix(model string, chat string) string {
	app := Dim + gptx.AppName + Reset
	model = Bold + G + model + Reset
	title := Bold + B + chat + Reset

	sep := Dim + "@" + Reset
	prefix := Dim + " ~> " + Reset
	postfix := Dim + " $ " + Reset

	if chat == "" {
		return fmt.Sprintf(
			Bold+"%s%s%s%s"+Reset, app, sep, model, postfix,
		)
	} else {
		return fmt.Sprintf(Bold+"%s%s%s%s%s%s"+Reset,
			app, sep, model, prefix, title, postfix,
		)
	}
}

// MARK: Logging
// ============================================================================

func doneMsg(msg string) {
	print(Bold + G) // style success prefix
	println(fmt.Sprintf("done: %s"+Reset, msg))
}

func errMsg(err error) {
	print(Bold + R) // style panic prefix
	println(fmt.Errorf("error: %w"+Reset, err).Error())
	print(Reset) // reset style
}

func warnMsg(err error) {
	print(Bold + Y) // style warning prefix
	println(fmt.Errorf("warning: %w"+Reset, err).Error())
}
