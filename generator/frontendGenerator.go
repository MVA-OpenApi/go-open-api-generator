package generator

import (
	fs "go-open-api-generator/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateFrontend(conf GeneratorConfig) {
	// create folder
	path := filepath.Join(conf.OutputPath, "public")
	fs.GenerateFolder(path)

	// create static html files
	createFileFromTemplate(filepath.Join(path, "index.html"), "templates/index.html.tmpl", conf)

	// copy OpenAPI Specification in this directory
	fs.CopyFile(conf.OpenAPIPath, path, conf.ModuleName + ".yaml")

	log.Info().Msg("Created Frontend.")
}