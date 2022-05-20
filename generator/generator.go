package generator

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"text/template"

	fs "go-open-api-generator/fileUtils"
	"go-open-api-generator/parser"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

const (
	Build       = "build"
	Cmd         = "cmd"
	Pkg         = "pkg"
	DefaultPort = 3000
)

type PortConfig struct {
	Port int16
}

func GenerateServer(openAPIFilePath string) {
	spec, err := parser.ParseOpenAPISpecFile(openAPIFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load OpenAPI spec file")
		return
	}

	createBuildDirectory()

	generateServerTemplate(spec.Servers[0].Variables["port"])
}

func createBuildDirectory() {
	// Removes previously generated folder structure
	fs.DeleteFolderRecursively(Build)

	// Generates basic folder structure
	fs.GenerateFolder(Build)
	fs.GenerateFolder(filepath.Join(Build, Cmd))
	fs.GenerateFolder(filepath.Join(Build, Pkg))
}

func generateServerTemplate(portSpec *openapi3.ServerVariable) {
	vars := PortConfig{DefaultPort}

	if portSpec != nil {
		portStr := portSpec.Default
		if portSpec.Enum != nil {
			portStr = portSpec.Enum[0]
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Warn().Msg("Failed to convert port, using 3000 instead")
		} else {
			vars.Port = int16(port)
		}
	}

	fileName := "main.go"
	templateFile := "templates/server.go.tmpl"
	templateName := path.Base(templateFile)

	// Create main.go and open it
	mainPath := filepath.Join(Build, Cmd, fileName)
	fs.GenerateFile(mainPath)
	file, fErr := os.OpenFile(mainPath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		fmt.Println(fErr.Error())
	}
	defer file.Close()

	// Parse the tempalte and write into main.go
	tmpl := template.Must(template.New(templateName).ParseFiles(templateFile))
	tmplErr := tmpl.Execute(file, vars)
	if tmplErr != nil {
		fmt.Println(tmplErr.Error())
	}
}
