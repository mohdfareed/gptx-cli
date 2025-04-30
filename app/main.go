package main

import (
	"fmt"
	"os"

	chat "github.com/mohdfareed/chatgpt-cli/pkg"
)

func main() {
	message := os.Args[len(os.Args)-1]
	fmt.Println(chat.Send(message))
}
