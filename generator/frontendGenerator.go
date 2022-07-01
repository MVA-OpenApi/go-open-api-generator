package generator

import (
	"fmt"
	fs "go-open-api-generator/fileUtils"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

func generateFrontend(spec *openapi3.T, conf GeneratorConfig) {
	generateOpenAPIDoc(conf)

	// create Schemas struct and add SchemaConfs (with name and properties) for schemas with x-label: form
	schemas := createSchemas(spec)

	schemasWithMethods := addMethods(schemas)

	fmt.Println(schemasWithMethods.List)

	// create folders
	frontendPath := filepath.Join(conf.OutputPath, "frontend")
	frontendPublicPath := filepath.Join(frontendPath, "public")
	frontendSrcPath := filepath.Join(frontendPath, "src")
	frontendStylesPath := filepath.Join(frontendSrcPath, "styles")
	frontendComponentsPath := filepath.Join(frontendSrcPath, "components")
	frontendSchemasPath := filepath.Join(frontendComponentsPath, "schemas")
	frontendSchemaFormsPath := filepath.Join(frontendSchemasPath, "schemaforms")

	fs.GenerateFolder(frontendPath)
	fs.GenerateFolder(frontendPublicPath)
	fs.GenerateFolder(frontendSrcPath)
	fs.GenerateFolder(frontendStylesPath)
	fs.GenerateFolder(frontendComponentsPath)
	fs.GenerateFolder(frontendSchemasPath)
	fs.GenerateFolder(frontendSchemaFormsPath)

	// create "static templates"
	createFileFromTemplate(filepath.Join(frontendPath, ".gitignore"), "templates/frontend/gitignore.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendPath, "package.json"), "templates/frontend/package.json.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/frontend/README.md.tmpl", nil)

	createFileFromTemplate(filepath.Join(frontendStylesPath, "index.css"), "templates/frontend/src/styles/index.css.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendStylesPath, "loadsubmit.css"), "templates/frontend/src/styles/loadsubmit.css.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendStylesPath, "schemaform.css"), "templates/frontend/src/styles/schemaform.css.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendStylesPath, "sidebar.css"), "templates/frontend/src/styles/sidebar.css.tmpl", nil)

	createFileFromTemplate(filepath.Join(frontendSrcPath, "api.js"), "templates/frontend/src/api.js.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendSrcPath, "index.js"), "templates/frontend/src/index.js.tmpl", nil)

	createFileFromTemplate(filepath.Join(frontendComponentsPath, "IDForm.js"), "templates/frontend/src/components/IDForm.js.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendComponentsPath, "LoadSubmit.js"), "templates/frontend/src/components/LoadSubmit.js.tmpl", nil)

	// create App.js
	createFileFromTemplate(filepath.Join(frontendSrcPath, "App.js"), "templates/frontend/src/App.js.tmpl", schemasWithMethods)

	// create SideBar.js
	createFileFromTemplate(filepath.Join(frontendComponentsPath, "SideBar.js"), "templates/frontend/src/components/SideBar.js.tmpl", schemasWithMethods)

	// create a schema componenta and a schema form components for each schema inside schemasWithMethods

	log.Info().Msg("Created Frontend.")
}

// for each schema in schemas.List add CRUD Methods with RESTful best practice API endpoints
func addMethods(schemas Schemas) (schemasWithMethods Schemas) {
	schemasWithMethods = schemas

	for i := range schemasWithMethods.List {
		tmpSchema := schemasWithMethods.List[i]
		tmpSchemaMethods := make([]MethodConf, 0)
		schemaURL := strings.ReplaceAll(strings.ToLower(tmpSchema.Name), " ", "")

		// add GET with path /<schema name>
		var getConf MethodConf
		getConf.Type = "get"
		getConf.Endpoint = schemaURL
		getConf.BodySchemaRequired = false
		tmpSchemaMethods = append(tmpSchemaMethods, getConf)

		// add GET with path /<schema name>/:id
		var getSpecificConf MethodConf
		getSpecificConf.Type = "get"
		getSpecificConf.Endpoint = schemaURL
		getSpecificConf.PathParams = make(map[string]string)
		getSpecificConf.PathParams["id"] = tmpSchema.Properties["id"]
		getSpecificConf.BodySchemaRequired = false
		tmpSchemaMethods = append(tmpSchemaMethods, getSpecificConf)

		// add POST with path /<schema name>
		var postConf MethodConf
		postConf.Type = "post"
		postConf.Endpoint = schemaURL
		postConf.BodySchemaRequired = true
		tmpSchemaMethods = append(tmpSchemaMethods, postConf)

		// add PUT with path /<schema name>/:id
		var putConf MethodConf
		putConf.Type = "put"
		putConf.Endpoint = schemaURL
		putConf.PathParams = make(map[string]string)
		putConf.PathParams["id"] = tmpSchema.Properties["id"]
		putConf.BodySchemaRequired = true
		tmpSchemaMethods = append(tmpSchemaMethods, putConf)

		// add DELETE with path /<schema name>/:id
		var deleteConf MethodConf
		deleteConf.Type = "delete"
		deleteConf.Endpoint = schemaURL
		deleteConf.PathParams = make(map[string]string)
		deleteConf.PathParams["id"] = tmpSchema.Properties["id"]
		deleteConf.BodySchemaRequired = false
		tmpSchemaMethods = append(tmpSchemaMethods, deleteConf)

		schemasWithMethods.List[i].Methods = tmpSchemaMethods

	}

	return schemasWithMethods

}

func createSchemas(spec *openapi3.T) (schemas Schemas) {
	schemas.List = make([]SchemaConf, 0)
	schemas.IsNotEmpty = false

	schemaStrings := toString(reflect.ValueOf(spec.Components.Schemas).MapKeys())

	for i := range schemaStrings {
		tmpSchemaName := schemaStrings[i]

		// check if schema has x-label == "form" -> if yes add schema to list
		schemaInformation, _ := spec.Components.Schemas[tmpSchemaName].Value.MarshalJSON()
		if strings.Contains(string(schemaInformation[:]), "\"x-label\":\"form\"") {
			var schema SchemaConf
			schema.Properties = make(map[string]string)

			// add names
			schema.Name = strings.ReplaceAll(strings.ToLower(tmpSchemaName), " ", "")
			schema.H1Name = strings.Title(tmpSchemaName)
			schema.ComponentName = strings.ReplaceAll(schema.H1Name, " ", "")

			// add properties
			tmpSchemaPropertyNames := reflect.ValueOf(spec.Components.Schemas[tmpSchemaName].Value.Properties).MapKeys()
			for j := range tmpSchemaPropertyNames {
				tmpSchemaPropertyName := tmpSchemaPropertyNames[j].Interface().(string)
				schema.Properties[strings.Title(tmpSchemaPropertyName)] = spec.Components.Schemas[tmpSchemaName].Value.Properties[tmpSchemaPropertyName].Value.Type
			}

			schemas.List = append(schemas.List, schema)
			schemas.IsNotEmpty = true
		}

	}

	return schemas

}

// function to convert an []reflect.Value to []string
func toString(inputArray []reflect.Value) (resultArray []string) {
	for i := range inputArray {
		resultArray = append(resultArray, inputArray[i].Interface().(string))
	}

	return resultArray
}

func generateOpenAPIDoc(conf GeneratorConfig) {
	// create folder
	type templateConfig struct {
		GeneratorConfig
		OpenAPIFile string
	}
	path := filepath.Join(conf.OutputPath, "public")
	fs.GenerateFolder(path)

	template := templateConfig{
		GeneratorConfig: conf,
		OpenAPIFile:     fs.GetFileNameWithEnding(conf.OpenAPIPath),
	}

	// create static html files
	createFileFromTemplate(filepath.Join(path, "index.html"), "templates/index.html.tmpl", template)

	// copy OpenAPI Specification in this directory
	fs.CopyFile(conf.OpenAPIPath, path, template.OpenAPIFile)

	log.Info().Msg("Created OpenAPI Documentation successfully.")
}
