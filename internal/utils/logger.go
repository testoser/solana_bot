package utils

import (
	"log"
	"os"
)

// NewLogger creates a new logger
func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	return logger
}
