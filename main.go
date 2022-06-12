package main

import (
	cli "go-open-api-generator/cli"
	generator "go-open-api-generator/generator"

	"embed"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed templates
var tmplFS embed.FS

func main() {
	// Set up zerolog time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Set pretty logging on
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// Export embed template filesystem to generator package
	generator.TmplFS = tmplFS

	cli.Execute()
}
