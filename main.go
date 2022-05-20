package main

import (
	cli "go-open-api-generator/cli"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set up zerolog time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Set pretty logging on
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cli.Execute()
}
