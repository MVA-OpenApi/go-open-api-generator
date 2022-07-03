package generator

import (
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	fs "go-open-api-generator/fileUtils"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func createFileFromTemplate(filePath string, tmplPath string, config interface{}) {
	templateName := path.Base(tmplPath)

	// Create file and open it
	fs.GenerateFile(filePath)
	file, fErr := os.OpenFile(filePath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		log.Fatal().Err(fErr).Msg("Failed creating file.")
		panic(fErr)
	}
	defer file.Close()

	// Parse the template and write into file
	tmpl := template.Must(template.New(templateName).ParseFS(TmplFS, tmplPath))
	tmplErr := tmpl.Execute(file, config)
	if tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing template.")
		panic(tmplErr)
	}

	log.Info().Msg("CREATE " + filePath)
}

func createFileFromTemplates(filePath string, tmplPaths []string, config interface{}) {
	templateName := path.Base(tmplPaths[0])

	// Create file and open it
	fs.GenerateFile(filePath)
	file, fErr := os.OpenFile(filePath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		log.Fatal().Err(fErr).Msg("Failed creating file.")
		panic(fErr)
	}
	defer file.Close()

	// Parse the template and write into file
	tmpl := template.Must(template.New(templateName).ParseFS(TmplFS, tmplPaths...))
	tmplErr := tmpl.Execute(file, config)
	if tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing template.")
		panic(tmplErr)
	}

	log.Info().Msg("CREATE " + filePath)
}

func validateStatusCode(code string) bool {
	return regexp.MustCompile(`[1-5](\d\d|XX)`).MatchString(code)
}

func convertPathParams(path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
}

func updateAuthConfig(spec *openapi3.T, conf *GeneratorConfig) {
	for key, value := range spec.Components.SecuritySchemes {
		if value.Value.Type == "apiKey" {
			conf.UseAuth = true
			conf.ApiKeyHeaderName = value.Value.Name
			conf.ApiKeySecurityName = key
			return
		}
	}
}

func updateOAPIOperation(op *openapi3.Operation, opID string, opSummary string, opDefault int) {
	op.OperationID = opID
	op.Summary = opSummary
	op.Responses.Default().Value = op.Responses.Get(opDefault).Value
}

func createOAPIResponse(rDesc string) *openapi3.Response {
	r := openapi3.NewResponse()
	r.Description = &rDesc
	return r
}
