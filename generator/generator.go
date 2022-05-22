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

	log.Info().Msg("Created all files successfully.")
}

func createBuildDirectory() {
	// Removes previously generated folder structure
	fs.DeleteFolderRecursively(Build)

	// Generates basic folder structure
	fs.GenerateFolder(Build)
	fs.GenerateFolder(filepath.Join(Build, Cmd))
	fs.GenerateFolder(filepath.Join(Build, Pkg))
	fs.GenerateFolder(filepath.Join(Build, Pkg, HandlerPkg))

	log.Info().Msg("Created project build directory.")
}

func generateServerTemplate(portSpec *openapi3.ServerVariable) {
	conf := PortConfig{DefaultPort}

	if portSpec != nil {
		portStr := portSpec.Default
		if portSpec.Enum != nil {
			portStr = portSpec.Enum[0]
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Warn().Msg("Failed to convert port, using 3000 instead.")
		} else {
			conf.Port = int16(port)
		}
	}

	fileName := "main.go"
	filePath := filepath.Join(Build, Cmd, fileName)
	templateFile := "templates/server.go.tmpl"

	log.Info().Msg("Creating server at port " + strconv.Itoa(int(conf.Port)) + "...")
	createFileFromTemplate(filePath, templateFile, conf)
}

func generateHandlerFuncStub(op *openapi3.Operation) OperationConfig {
	var conf OperationConfig

	conf.Summary = op.Summary
	conf.OperationID = op.OperationID

	for resKey, resRef := range op.Responses {
		conf.Responses = append(conf.Responses, ResponseConfig{resKey, *resRef.Value.Description})
	}

	fileName := conf.OperationID + ".go"
	filePath := filepath.Join(Build, Pkg, HandlerPkg, fileName)
	templateFile := "templates/handlerFunc.go.tmpl"

	createFileFromTemplate(filePath, templateFile, conf)

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
	filePath := filepath.Join(Build, Pkg, HandlerPkg, fileName)
	templateFile := "templates/handler.go.tmpl"

	createFileFromTemplate(filePath, templateFile, conf)
}

func createFileFromTemplate(filePath string, tmplPath string, config any) {
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
