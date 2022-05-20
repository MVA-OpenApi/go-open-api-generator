package parser

import (
	"errors"
	fs "go-open-api-generator/fileUtils"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func ParseOpenAPISpecFile(path string) (*openapi3.T, error) {
	if !fs.CheckIfFileExists(path) {
		return nil, errors.New("file not found")
	}

	loader := openapi3.NewLoader()

	spec, err := loader.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("OpenAPI Spec file loaded successfully")

	err = spec.Validate(loader.Context)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("OpenAPI Spec file validated successfully")

	return spec, err
}
