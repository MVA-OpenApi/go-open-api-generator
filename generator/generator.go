package generator

import (
	"errors"
	"path/filepath"
	"strconv"

	fs "go-open-api-generator/fileUtils"
	"go-open-api-generator/parser"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

const (
	Cmd         = "cmd"
	Pkg         = "pkg"
	UtilPkg     = "util"
	HandlerPkg  = "handler"
	DatabasePkg = "db"
	ModelPkg    = "model"
	DefaultPort = 3000
)

var (
	config ProjectConfig
)

func GenerateServer(conf GeneratorConfig) error {
	spec, err := parser.ParseOpenAPISpecFile(conf.OpenAPIPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load OpenAPI spec file")
		return err
	}

	// Init project config
	config.Name = conf.ModuleName
	config.Path = conf.OutputPath

	createProjectPathDirectory()

	serverConf := generateServerTemplate(spec.Servers[0].Variables["port"], conf)

	generateConfigFiles(serverConf)

	generateFrontend(conf)

	generateHandlerFuncs(spec)

	GenerateTypes(spec, config)

	generateDatabaseFiles(conf)

	log.Info().Msg("Created all files successfully.")

	return nil
}

func createProjectPathDirectory() {
	// Removes previously generated folder structure
	fs.DeleteFolderRecursively(config.Path)

	// Generates basic folder structure
	fs.GenerateFolder(config.Path)
	fs.GenerateFolder(filepath.Join(config.Path, Pkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, UtilPkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, HandlerPkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, DatabasePkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, ModelPkg))

	log.Info().Msg("Created project directory.")
}

func generateServerTemplate(portSpec *openapi3.ServerVariable, generatorConf GeneratorConfig) (serverConf ServerConfig) {
	openAPIName := fs.GetFileName(generatorConf.OpenAPIPath)
	conf := ServerConfig{
		Port:        DefaultPort,
		ModuleName:  config.Name,
		Flags:       generatorConf.Flags,
		OpenAPIName: openAPIName,
	}

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
	} else {
		log.Warn().Msg("No port field was found, using 3000 instead.")
	}

	if generatorConf.UseLogger {
		log.Info().Msg("Adding logging middleware.")
	}

	fileName := "main.go"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/server.go.tmpl"

	log.Info().Msg("Creating server at port " + strconv.Itoa(int(conf.Port)) + "...")
	createFileFromTemplate(filePath, templateFile, conf)

	return conf
}

func generateHandlerFuncStub(op *openapi3.Operation, method string, path string) (OperationConfig, error) {
	var conf OperationConfig
	var methodPath = method + " " + path

	conf.Method = method

	conf.Summary = op.Summary
	if op.Summary == "" {
		log.Warn().Msg("No summary found for endpoint: " + methodPath)
	}

	conf.OperationID = op.OperationID
	if op.OperationID == "" {
		log.Error().Msg("No operation ID found for endpoint: " + methodPath)
		return conf, errors.New("no operation id, can't create function")
	}

	for resKey, resRef := range op.Responses {
		if !validateStatusCode(resKey) {
			log.Warn().Msg("Status code " + resKey + " for endpoint " + methodPath + " is not a valid status code.")
		}

		conf.Responses = append(conf.Responses, ResponseConfig{resKey, *resRef.Value.Description})
	}

	fileName := conf.OperationID + ".go"
	filePath := filepath.Join(config.Path, Pkg, HandlerPkg, fileName)
	templateFile := "templates/handlerFunc.go.tmpl"

	createFileFromTemplate(filePath, templateFile, conf)

	return conf, nil
}

func generateHandlerFuncs(spec *openapi3.T) {
	var conf HandlerConfig
	for path, pathObj := range spec.Paths {
		var newPath PathConfig
		newPath.Path = convertPathParams(path)

		for method, op := range pathObj.Operations() {
			opConfig, err := generateHandlerFuncStub(op, method, newPath.Path)

			if err != nil {
				log.Err(err).Msg("Skipping generation of handler function for endpoint " + method + " " + path)
			}

			newPath.Operations = append(newPath.Operations, opConfig)
		}

		conf.Paths = append(conf.Paths, newPath)
	}

	fileName := "handler.go"
	filePath := filepath.Join(config.Path, Pkg, HandlerPkg, fileName)
	templateFile := "templates/handler.go.tmpl"

	createFileFromTemplate(filePath, templateFile, conf)
}

func generateConfigFiles(serverConf ServerConfig) {
	// create app.env file
	fileName := "app.env"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/app.env.tmpl"

	createFileFromTemplate(filePath, templateFile, serverConf)

	// create config.go.tmpl file
	fileName = "config.go"
	filePath = filepath.Join(config.Path, Pkg, UtilPkg, fileName)
	templateFile = "templates/config.go.tmpl"

	createFileFromTemplate(filePath, templateFile, nil)

}

func generateDatabaseFiles(conf GeneratorConfig) {
	if conf.UseDatabase {
		log.Info().Msg("Adding SQLite database.")
	}

	fileName := conf.DatabaseName
	filePath := filepath.Join(config.Path, Pkg, DatabasePkg, fileName)
	templateFile := "templates/database.go.tmpl"

	fs.GenerateFile(filePath + ".db")
	createFileFromTemplate(filePath+".go", templateFile, conf)
}
