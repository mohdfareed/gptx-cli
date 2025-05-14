//go:build dev
// +build dev

package cfg

import (
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}
