package generator

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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
	HandlerPkg  = "handler"
	DefaultPort = 3000
)

func GenerateServer(openAPIFilePath string) {
	spec, err := parser.ParseOpenAPISpecFile(openAPIFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load OpenAPI spec file")
		return
	}

	createBuildDirectory()

	generateServerTemplate(spec.Servers[0].Variables["port"])

	generateHandlerFuncs(spec)
}

func createBuildDirectory() {
	// Removes previously generated folder structure
	fs.DeleteFolderRecursively(Build)

	// Generates basic folder structure
	fs.GenerateFolder(Build)
	fs.GenerateFolder(filepath.Join(Build, Cmd))
	fs.GenerateFolder(filepath.Join(Build, Pkg))
	fs.GenerateFolder(filepath.Join(Build, Pkg, HandlerPkg))
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
		log.Fatal().Err(fErr).Msg("Failed creating main.go file")
		panic(fErr)
	}
	defer file.Close()

	// Parse the template and write into main.go
	tmpl := template.Must(template.New(templateName).ParseFiles(templateFile))
	tmplErr := tmpl.Execute(file, vars)
	if tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing template")
		panic(tmplErr)
	}
}

func generateHandlerFuncStub(op *openapi3.Operation) OperationConfig {
	var conf OperationConfig

	conf.Summary = op.Summary
	conf.OperationID = op.OperationID

	for _, responseRef := range op.Responses {
		conf.Responses = append(conf.Responses, ResponseConfig{responseRef.Ref, *responseRef.Value.Description})
	}

	fileName := conf.OperationID + ".go"
	templateFile := "templates/handlerFunc.go.tmpl"
	templateName := path.Base(templateFile)

	// Create handler func file and open it
	mainPath := filepath.Join(Build, Pkg, HandlerPkg, fileName)
	fs.GenerateFile(mainPath)
	file, fErr := os.OpenFile(mainPath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		log.Fatal().Err(fErr).Msg("Failed creating file")
		panic(fErr)
	}
	defer file.Close()

	// Parse the template and write into main.go
	tmpl := template.Must(template.New(templateName).ParseFiles(templateFile))
	tmplErr := tmpl.Execute(file, conf)
	if tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing template")
		panic(tmplErr)
	}

	return conf
}

func generateHandlerFuncs(spec *openapi3.T) {
	var conf HandlerConfig

	for path, pathObj := range spec.Paths {
		var newPath PathConfig
		newPath.Path = strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")

		for method, op := range pathObj.Operations() {
			opConfig := generateHandlerFuncStub(op)
			opConfig.Method = method

			newPath.Operations = append(newPath.Operations, opConfig)
		}

		conf.Paths = append(conf.Paths, newPath)
	}

	fileName := "handler.go"
	templateFile := "templates/handler.go.tmpl"
	templateName := path.Base(templateFile)

	// Create handler.go and open it
	mainPath := filepath.Join(Build, Pkg, HandlerPkg, fileName)
	fs.GenerateFile(mainPath)
	file, fErr := os.OpenFile(mainPath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		log.Fatal().Err(fErr).Msg("Failed creating file")
		panic(fErr)
	}
	defer file.Close()

	// Parse the template and write into main.go
	tmpl := template.Must(template.New(templateName).ParseFiles(templateFile))
	tmplErr := tmpl.Execute(file, conf)
	if tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing template")
		panic(tmplErr)
	}
}
