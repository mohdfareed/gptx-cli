//go:build dev
// +build dev

package gptx

import (
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}
