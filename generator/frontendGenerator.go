package generator

import (
	fs "go-open-api-generator/fileUtils"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func generateFrontend(conf GeneratorConfig) {
	// create folder
	type templateConfig struct {
		GeneratorConfig
		OpenAPIName string
	}
	path := filepath.Join(conf.OutputPath, "public")
	fs.GenerateFolder(path)

	template := templateConfig{
		GeneratorConfig: conf,
		OpenAPIName:     fs.GetFileName(conf.OpenAPIPath),
	}

	// create static html files
	createFileFromTemplate(filepath.Join(path, "index.html"), "templates/index.html.tmpl", template)

	// copy OpenAPI Specification in this directory
	fs.CopyFile(conf.OpenAPIPath, path, template.OpenAPIName+".yaml")

	log.Info().Msg("Created Frontend.")
}
