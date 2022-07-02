package generator

import (
	"bufio"
	"embed"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
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
	Cmd                      = "cmd"
	Pkg                      = "pkg"
	UtilPkg                  = "util"
	HandlerPkg               = "handler"
	DatabasePkg              = "db"
	ModelPkg                 = "model"
	AuthzPkg                 = "authz"
	MiddlewarePackage        = "middleware"
	DefaultPort              = 8080
	PUT               string = "\"PUT\""
	GET                      = "\"GET\""
	POST                     = "\"POST\""
	DELETE                   = "\"DELETE\""
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

	generateBdd("../tests/stores.feature")

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

//---------GENERATION OF BDD---------

func ignore(input string) bool {
	return input == "When" || input == "And" || input == "Given" || input == "Then"
}

func retrieveRegex(input string) string {
	regex := ""
	for _, j := range input {
		if string(j) == "{" {
			regex = regex + "\\\\" + string(j)
		} else {
			regex = regex + string(j)
		}
	}
	return regex
}

func parseSteps(path string) []Step {
	//We use this map to connect each Step struct which represents a step in godog, to the answer that it requires
	m := make(map[string]int)
	var listOfSteps []Step
	file, err := os.Open(path)
	if err != nil {
		log.Fatal()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var prevStepConf Step
	regexedPath := ""

	for scanner.Scan() {
		var stepConf Step
		stepConf.Mapping = make(map[string]int)
		stepConf.RegexAndCode = make(map[string]int)
		line := scanner.Text()
		words := strings.Fields(line)
		stringRegex := "\"([^\"]*)\""
		for i, word := range words {
			//We skip this since we don't need to create any method for them
			if word == "Scenario:" && i == 0 {
				for _, j := range words[i:] {
					if ok, _ := regexp.MatchString(stringRegex, j); ok {
						regexedPath = retrieveRegex(j)
					}
				}
				break
			}
			if word == "Feature:" && i == 0 {
				break
			} else {
				//Retrieve the Method being used
				if word == PUT || word == GET || word == POST || word == DELETE {
					stepConf.Name = stepConf.Name + word
					word = strings.ToLower(word)
					r := []rune(word)
					stepConf.Method = string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
				} else if i >= 1 && (words[i-1] == "to") {
					//After "to" in the predefined structure we always receive the endpoint
					for _, word := range words[i:] {
						if value, _ := regexp.MatchString(stringRegex, word); value {
							stepConf.Endpoint = stepConf.Endpoint + word
							break
						}
					}
				} else if n, err := strconv.Atoi(word); err == nil && 200 <= n && n <= 600 {
					//Retrieve the status code
					stepConf.StatusCode = word
				} else if i >= 1 && (words[i-1] == "payload" || words[i-1] == "Payload" || words[i-1] == "PAYLOAD") {
					//After the word "payload" in the predefined structure comes the payload of the request
					for _, j := range words[i:] {
						stepConf.Payload = stepConf.Payload + j
					}
					//If we are in a "Then" part of the scenario, the next word is "the", there receive the status code
					//We also change the values in the previous step. We do this because the previous step, the actual
					//method is where we send the requests to the server whereas the "Then" part of the scenario
					//is just used to check whether the response was correct based on the request's method, url and payload
					if i == len(words)-1 && len(stepConf.StatusCode) == 0 {
						prevStepConf = stepConf
					}
					break
				} else {
					//Fill the name of the step. ignore() is used to ignore keywords such as When, Then, Given, And, since
					//they are not part of the name
					ignore := ignore(word)
					if len(stepConf.Name) == 0 && !ignore {
						//if the name is empty we add the word in it but firstly it has to be lower cased since the name uses
						//camel case convention
						stepConf.Name = stepConf.Name + strings.ToLower(word)
					} else if len(stepConf.Name) != 0 && !ignore {
						//if the name is not empty we capitalize the first letter of the word and add the word to the name
						//we also want to make sure that the name does not contain any substrings in quotes
						value, _ := regexp.MatchString(stringRegex, word)
						if !value {
							r := []rune(word)
							stepConf.Name = stepConf.Name + string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
						}
					}
				}
				//we receive the status code of the step in the "Then" part of the scenario. That's why we use prevStepConf.
				if strings.HasPrefix(stepConf.Name, "the") {
					code, _ := strconv.Atoi(stepConf.StatusCode)
					if code != 0 && len(regexedPath) != 0 {
						m[prevStepConf.Endpoint] = code
						prevStepConf.RegexPaths = append(prevStepConf.RegexPaths, regexedPath)
						prevStepConf.StatusCode = strconv.Itoa(code)
						value, _ := strconv.Atoi(prevStepConf.StatusCode)
						prevStepConf.Mapping[prevStepConf.Endpoint] = value
						prevStepConf.RegexAndCode[regexedPath] = value
					}
				}
				//Finally if we have arrived at the end of the string we are looking to add prevStepConf in our list of steps
				if i == len(words)-1 {
					flag := false
					code, _ := strconv.Atoi(stepConf.StatusCode)
					//Because the function handlers are not implemented, the payloads do not play an important role. Endpoints do.
					m[stepConf.Endpoint] = code
					if strings.HasPrefix(stepConf.Name, "the") && len(prevStepConf.Name) != 0 {
						//if the list of steps is empty then we just add the step
						if len(listOfSteps) == 0 {
							listOfSteps = append(listOfSteps, prevStepConf)
						} else {
							//otherwise we check whether there is a step with the same name or not
							for _, j := range listOfSteps {
								//if there is, then we want to add this step to the mapping of the step with the same name
								//thus we don't add redundant steps in the list, and also
								//it is exactly like the schema that godog requires in order to work
								if prevStepConf.Name == j.Name {
									flag = true
									value, _ := strconv.Atoi(prevStepConf.StatusCode)
									j.Mapping[prevStepConf.Endpoint] = value
									j.RegexAndCode[regexedPath] = value
									j.RegexPaths = append(j.RegexPaths, prevStepConf.RegexPaths[0])
								}
							}
							//in the case where the name is not there and also the list is not empty we use the flag to check
							//and we just append the step to the list
							if flag == false {
								listOfSteps = append(listOfSteps, prevStepConf)
							}
						}
					}
				}
			}
			//after we're done with the initialization of the stepConf, we initialize prevStepConf and move on
			if i == len(words)-1 && len(stepConf.StatusCode) == 0 {
				prevStepConf = stepConf
			}
		}
	}
	//add the string that we want to use in the godog file
	for i, v := range listOfSteps {
		v.RealName = v.RealName + createName(v.Name)
		localhost := "http://localhost:8080"
		for _, k := range v.RegexPaths {
			in := strings.ReplaceAll(localhost+k, "\"", "")
			v.PathsWithHost = append(v.PathsWithHost, in)
		}
		listOfSteps[i] = v
	}
	return listOfSteps
}

//add space between words in camelcase
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func AddedSpace(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1} ${2}")
	return strings.ToLower(snake)
}

//we use this to create the string for the function which will be used in InitializeScenario() function
func createName(str string) string {
	nameToReturn := ""
	transformed := AddedSpace(str)
	for i, word := range strings.Fields(transformed) {
		if word == "i" && i == 0 {
			nameToReturn = nameToReturn + "^" + strings.ToUpper(word) + " "
		} else if word == "to" {
			nameToReturn = nameToReturn + "to \"([^\"]*)\" "
		} else if word == "payload" {
			nameToReturn = nameToReturn + "payload \"([^\"]*)\""
		} else if i == len(transformed)-1 {
			nameToReturn = nameToReturn + "$"
		} else if word == "put" || word == "get" || word == "post" || word == "delete" {
			nameToReturn = nameToReturn + strings.ToUpper(word) + " "
		} else {
			nameToReturn = nameToReturn + word + " "
		}
	}
	return nameToReturn
}

func contains(element string, arr []string) bool {
	for _, j := range arr {
		if element == j {
			return true
		}
	}
	return false
}

func getAllEndpoints(listing Listing) []string {
	var slice []string
	for _, k := range listing.Steps {
		for _, j := range k.RegexPaths {
			if !contains(j, slice) {
				slice = append(slice, j)
			}
		}
	}
	return slice
}

func generateBdd(path string) {
	var step Listing
	step.Steps = parseSteps(path)
	step.UniqueEndpoints = getAllEndpoints(step)

	f, err := os.OpenFile("generation_godog_test.go", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	content, _ := ioutil.ReadFile("templates/bdd.go.tmpl")
	t := template.Must(template.New("bdd-tmpl").Parse(string(content)))
	err1 := t.Execute(f, step)
	if err1 != nil {
		panic(err1)
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
