package generator

import (
	"github.com/rs/zerolog/log"
	filepath2 "path/filepath"
)

func generateTemplates(name string, configuration GeneratorConfig) {
	templateName := "templates/bdd.go.tmpl"
	featureName := name + "_feature.feature"
	filepath := filepath2.Join(config.Path, featureName)
	log.Info().Msg("Creating feature file at path" + filepath)
	createFileFromTemplate(filepath, templateName, configuration)
}
