package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr)
	logger.Info().Msg("Hello there!")
}
