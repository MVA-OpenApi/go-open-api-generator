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
	Cmd         = "cmd"
	Pkg         = "pkg"
	HandlerPkg  = "handler"
	DefaultPort = 3000
)

var (
	config ProjectConfig
)

func GenerateServer(openAPIPath string, projectPath string, moduleName string) {
	spec, err := parser.ParseOpenAPISpecFile(openAPIPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load OpenAPI spec file")
		return
	}

	// Init project config
	config.Name = moduleName
	config.Path = projectPath

	createProjectPathDirectory()

	generateServerTemplate(spec.Servers[0].Variables["port"])

	generateHandlerFuncs(spec)

	GenerateTypes(spec, config)

	log.Info().Msg("Created all files successfully.")
}

func createProjectPathDirectory() {
	// Removes previously generated folder structure
	fs.DeleteFolderRecursively(config.Path)

	// Generates basic folder structure
	fs.GenerateFolder(config.Path)
	fs.GenerateFolder(filepath.Join(config.Path, Cmd))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, HandlerPkg))

	log.Info().Msg("Created project directory.")
}

func generateServerTemplate(portSpec *openapi3.ServerVariable) {
	conf := ServerConfig{Port: DefaultPort, ModuleName: config.Name}

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
	filePath := filepath.Join(config.Path, Cmd, fileName)
	templateFile := "templates/server.go.tmpl"

	log.Info().Msg("Creating server at port " + strconv.Itoa(int(conf.Port)) + "...")
	CreateFileFromTemplate(filePath, templateFile, conf)
}

func generateHandlerFuncStub(op *openapi3.Operation) OperationConfig {
	var conf OperationConfig

	conf.Summary = op.Summary
	conf.OperationID = op.OperationID

	for resKey, resRef := range op.Responses {
		conf.Responses = append(conf.Responses, ResponseConfig{resKey, *resRef.Value.Description})
	}

	fileName := conf.OperationID + ".go"
	filePath := filepath.Join(config.Path, Pkg, HandlerPkg, fileName)
	templateFile := "templates/handlerFunc.go.tmpl"

	CreateFileFromTemplate(filePath, templateFile, conf)

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
	filePath := filepath.Join(config.Path, Pkg, HandlerPkg, fileName)
	templateFile := "templates/handler.go.tmpl"

	CreateFileFromTemplate(filePath, templateFile, conf)
}

func CreateFileFromTemplate(filePath string, tmplPath string, config interface{}) {
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
