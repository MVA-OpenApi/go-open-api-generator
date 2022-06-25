package generator

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	//"go-open-api-generator/generator"
	"path/filepath"
	"strconv"

	fs "go-open-api-generator/fileUtils"
	"go-open-api-generator/parser"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

const (
	Cmd                = "cmd"
	Pkg                = "pkg"
	UtilPkg            = "util"
	HandlerPkg         = "handler"
	DatabasePkg        = "db"
	ModelPkg           = "model"
	AuthzPkg           = "authz"
	DefaultPort        = 8080
	PUT         string = "\"PUT\""
	GET                = "\"GET\""
	POST               = "\"POST\""
	DELETE             = "\"DELETE\""
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

	serverConf := generateServerTemplate(spec.Servers[0].Variables["port"], conf)

	generateConfigFiles(serverConf)

	generateBdd("../tests/stores.feature")

	generateFrontend(conf)

	generateHandlerFuncs(spec, conf)

	GenerateTypes(spec, config)

	if conf.UseDatabase {
		generateDatabaseFiles(conf)
	}

	if conf.UseAuth {
		generateAuthzFile(conf)
	}

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
	if conf.UseDatabase {
		fs.GenerateFolder(filepath.Join(config.Path, Pkg, DatabasePkg))
	}
	fs.GenerateFolder(filepath.Join(config.Path, Pkg, ModelPkg))
	if conf.UseAuth {
		fs.GenerateFolder(filepath.Join(config.Path, Pkg, AuthzPkg))
	}

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
			log.Warn().Msg("Failed to convert port, using 8080 instead.")
		} else {
			conf.Port = int16(port)
		}
	} else {
		log.Warn().Msg("No port field was found, using 8080 instead.")
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

func generateHandlerFuncs(spec *openapi3.T, genConf GeneratorConfig) {
	type handlerConf struct {
		HandlerConfig
		UseAuth    bool
		ModuleName string
	}
	var conf handlerConf
	conf.ModuleName = genConf.ModuleName
	conf.UseAuth = genConf.UseAuth

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

//---------GENERATION OF BDD---------

func matchString(pattern string, s string) (bool, error) {
	return regexp.MatchString(s, pattern)
}

//We don't want to have function names of form whenISendGETRequest() but rather iSendGETRequest

func ignore(input string) bool {
	return input == "When" || input == "And" || input == "Given" || input == "Then"
}

func generateBdd(path string) {
	//We use this map to connect each Step struct which represents a step in godog, to the answer that it requires
	file, err := os.Open(path)
	if err != nil {
		log.Fatal()
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var stepConf Step
		line := scanner.Text()
		words := strings.Fields(line)
		stringRegex := "\"([^\"]*)\""
		for i, word := range words {
			j := 0
			if (word == "Scenario:" && i == 0) || (word == "Feature:" && i == 0) {
				break
			} else if word == "When" || word == "And" {
				for _, word := range words {
					value, _ := matchString(stringRegex, word)
					ignore := ignore(word)
					if !value && !ignore {
						stepConf.Name = stepConf.Name + word
					} else {
						j = j + 1
						//j serves as a counter for how many arguments we get
						argument := "arg" + strconv.Itoa(j)
						stepConf.Arguments = append(stepConf.Arguments, argument)
					}
				}
			} else {
				if word == PUT || word == GET || word == POST || word == DELETE {
					stepConf.Method = word
				} else if i >= 1 && (words[i-1] == "url" || words[i-1] == "URL" || words[i-1] == "endpoint" || words[i-1] == "Endpoint") {
					value, err := matchString(stringRegex, word)
					if err != nil {
						return
					}
					if value {
						stepConf.Endpoint = word
					}
				} else if n, err := strconv.Atoi(word); err == nil && 200 <= n && n <= 500 {
					stepConf.StatusCode = word
				} else if i >= 1 && (words[i-1] == "payload" || words[i-1] == "Payload" || words[i-1] == "PAYLOAD") {
					value, err := matchString(stringRegex, word)
					if err != nil {
						return
					}
					if value {
						stepConf.Payload = word
					}
				} else {
					ignore := ignore(word)
					if stepConf.Name == "" && !ignore {
						stepConf.Name = stepConf.Name + strings.ToLower(word)
					} else {
						value, _ := matchString(stringRegex, word)
						if !value {
							r := []rune(word)
							stepConf.Name = stepConf.Name + string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
						}
					}
				}
			}
			j = 0
		}
		fmt.Println("Name: ", stepConf.Name)
		fmt.Println("Endpoint: ", stepConf.Endpoint)
		fmt.Println("Payload: ", stepConf.Payload)
		fmt.Println("Status Code: ", stepConf.StatusCode)
		fmt.Println("Method: ", stepConf.Method)
	}
}

//---------END---------

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
