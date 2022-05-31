package generator

import (
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	fs "go-open-api-generator/fileUtils"

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
	tmpl := template.Must(template.New(templateName).ParseFiles(tmplPath))
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
