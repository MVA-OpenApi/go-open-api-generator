package generator

import (
	"embed"
	"errors"
	"path/filepath"
	"strconv"

	fs "go-open-api-generator/fileUtils"
	"go-open-api-generator/parser"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

const (
	Cmd               = "cmd"
	Pkg               = "pkg"
	UtilPkg           = "util"
	HandlerPkg        = "handler"
	DatabasePkg       = "db"
	ModelPkg          = "model"
	AuthzPkg          = "authz"
	MiddlewarePackage = "middleware"
	DefaultPort       = 8080
)

var (
	config ProjectConfig
	TmplFS embed.FS
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

	updateAuthConfig(spec, &conf)

	createProjectPathDirectory(conf)

	serverConf := generateServerTemplate(spec, conf)

	generateConfigFiles(serverConf)

	generateFrontend(conf)

	if conf.UseLifecycle {
		generateLifecycleFiles(spec)
	}

	generateHandlerFuncs(spec, conf)

	GenerateTypes(spec, config)

	if conf.UseDatabase {
		generateDatabaseFiles(conf)
	}

	if conf.UseAuth {
		generateAuthzFile(conf)
	}

	if conf.UseValidation {
		generateValidation(conf)
	}

	generateMakefile(conf, serverConf)

	generateDockerfile(conf, serverConf)

	log.Info().Msg("Created all files successfully.")

	return nil
}

func createProjectPathDirectory(conf GeneratorConfig) {
	// Removes previously generated folder structure
	fs.DeleteFolderRecursively(config.Path)

	// Generates basic folder structure
	fs.GenerateFolder(config.Path)
	fs.GenerateFolder(filepath.Join(config.Path, Pkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, UtilPkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, HandlerPkg))
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, ModelPkg))
	if conf.UseDatabase {
		fs.GenerateFolder(filepath.Join(config.Path, Pkg, DatabasePkg))
	}
	if conf.UseAuth {
		fs.GenerateFolder(filepath.Join(config.Path, Pkg, AuthzPkg))
	}
	if conf.UseValidation {
		fs.GenerateFolder(filepath.Join(config.Path, Pkg, MiddlewarePackage))
	}

	log.Info().Msg("Created project directory.")
}

func generateServerTemplate(spec *openapi3.T, generatorConf GeneratorConfig) (serverConf ServerConfig) {
	openAPIName := fs.GetFileName(generatorConf.OpenAPIPath)
	conf := ServerConfig{
		Port:        DefaultPort,
		ModuleName:  config.Name,
		Flags:       generatorConf.Flags,
		OpenAPIName: openAPIName,
	}

	strDefaultPort := strconv.Itoa(DefaultPort)

	if spec.Servers != nil {
		serverSpec := spec.Servers[0]
		if portSpec := serverSpec.Variables["port"]; portSpec != nil {
			portStr := portSpec.Default
			if portSpec.Enum != nil {
				portStr = portSpec.Enum[0]
			}

			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Warn().Msg("Failed to convert port, using" + strDefaultPort + "instead.")
			} else {
				conf.Port = int16(port)
			}
		} else {
			log.Warn().Msg("No port field was found, using" + strDefaultPort + "instead.")
		}
	} else {
		log.Warn().Msg("No servers field was found, using" + strDefaultPort + "instead.")
	}

	if generatorConf.UseLogger {
		log.Info().Msg("Adding logging middleware.")
	}

	if generatorConf.UseHTTP2 {
		log.Info().Msg("Using HTTP/2 as default protocol")
	}

	fileName := "main.go"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/server.go.tmpl"

	log.Info().Msg("Creating server at port " + strconv.Itoa(int(conf.Port)) + "...")
	createFileFromTemplate(filePath, templateFile, conf)

	return conf
}

func generateHandlerFuncStub(op *openapi3.Operation, method string, path string, apiSecurityName string) (OperationConfig, error) {
	var conf OperationConfig
	var methodPath = method + " " + path

	if op.Security != nil {
		for _, item := range *op.Security {
			for key := range item {
				if key == apiSecurityName {
					conf.UseAuth = true
					break
				}
			}
		}
	}

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
		if !validateStatusCode(resKey) && resKey != "default" {
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

func generateHandlerFuncs(spec *openapi3.T, genConf GeneratorConfig) {
	conf := HandlerConfig{
		ModuleName:  genConf.ModuleName,
		OpenAPIPath: fs.GetFileNameWithEnding(genConf.OpenAPIPath),
		UseAuth:     genConf.UseAuth,
		Flags:       genConf.Flags,
	}
	conf.ModuleName = genConf.ModuleName
	conf.Flags = genConf.Flags

	for _, item := range spec.Security {
		for key := range item {
			if key == genConf.ApiKeySecurityName {
				conf.UseGlobalAuth = true
				break
			}
		}
	}

	for path, pathObj := range spec.Paths {
		var newPath PathConfig
		newPath.Path = convertPathParams(path)

		for method, op := range pathObj.Operations() {
			opConfig, err := generateHandlerFuncStub(op, method, newPath.Path, genConf.ApiKeySecurityName)

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
	fileName := ".env"
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
	log.Info().Msg("Adding SQLite database.")

	fileName := conf.DatabaseName
	filePath := filepath.Join(config.Path, Pkg, DatabasePkg, fileName)
	templateFile := "templates/database.go.tmpl"

	fs.GenerateFile(filePath + ".db")
	createFileFromTemplate(filePath+".go", templateFile, conf)
}

func generateAuthzFile(conf GeneratorConfig) {
	log.Info().Msg("Adding auth middleware.")

	fileName := "authz.go"
	filePath := filepath.Join(config.Path, Pkg, AuthzPkg, fileName)
	templateFile := "templates/authz.go.tmpl"

	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, conf)
}

func generateValidation(conf GeneratorConfig) {
	log.Info().Msg("Adding validation middleware.")

	fileName := "validation.go"
	filePath := filepath.Join(config.Path, Pkg, MiddlewarePackage, fileName)
	templateFile := "templates/validation.go.tmpl"

	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, conf)
}

func generateLifecycleFiles(spec *openapi3.T) {
	if spec.Paths.Find("/livez") == nil {
		log.Info().Msg("Generating default /livez endpoint.")

		op := openapi3.NewOperation()
		op.AddResponse(200, createOAPIResponse("The server is alive"))
		updateOAPIOperation(op, "getHealth", "Returns health-state of the server", 200)
		spec.AddOperation("/livez", "GET", op)
	}
	if spec.Paths.Find("/readyz") == nil {
		log.Info().Msg("Generating default /readyz endpoint.")

		op := openapi3.NewOperation()
		op.AddResponse(200, createOAPIResponse("The Service is ready"))
		op.AddResponse(500, createOAPIResponse("The Service is not ready"))
		updateOAPIOperation(op, "getReady", "Returns ready-state of the server", 200)
		spec.AddOperation("/readyz", "GET", op)
	}
}

func generateMakefile(conf GeneratorConfig, serverConf ServerConfig) {
	type makefileConfig struct {
		ModuleName string
		Port       int16
	}

	var makefileConf makefileConfig
	makefileConf.ModuleName = conf.ModuleName
	makefileConf.Port = serverConf.Port

	log.Info().Msg("Adding Makefile.")

	fileName := "Makefile"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/make-file.tmpl"

	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, makefileConf)
}

func generateDockerfile(conf GeneratorConfig, serverConf ServerConfig) {
	type dockerfileConfig struct {
		ModuleName string
		Port       int16
	}

	var dockerfileConf dockerfileConfig
	dockerfileConf.ModuleName = conf.ModuleName
	dockerfileConf.Port = serverConf.Port

	log.Info().Msg("Adding Dockerfile.")

	fileName := "Dockerfile"
	filePath := filepath.Join(config.Path, fileName)
	templateFile := "templates/docker-file.tmpl"

	fs.GenerateFile(filePath)
	createFileFromTemplate(filePath, templateFile, dockerfileConf)
}
