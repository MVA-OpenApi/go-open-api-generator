package generator

import (
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

	// files in root directory
	createFileFromTemplate(filepath.Join(frontendPath, ".gitignore"), "templates/frontend/gitignore.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendPath, "package.json"), "templates/frontend/package.json.tmpl", conf)
	createFileFromTemplate(filepath.Join(frontendPath, "README.md"), "templates/frontend/README.md.tmpl", nil)

	// files in styles directory
	createFileFromTemplate(filepath.Join(frontendStylesPath, "index.css"), "templates/frontend/src/styles/index.css.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendStylesPath, "loadsubmit.css"), "templates/frontend/src/styles/loadsubmit.css.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendStylesPath, "schemaform.css"), "templates/frontend/src/styles/schemaform.css.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendStylesPath, "sidebar.css"), "templates/frontend/src/styles/sidebar.css.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendStylesPath, "defaultpage.css"), "templates/frontend/src/styles/defaultpage.css.tmpl", nil)

	// files in src directory
	createFileFromTemplate(filepath.Join(frontendSrcPath, "api.js"), "templates/frontend/src/api.js.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendSrcPath, "index.js"), "templates/frontend/src/index.js.tmpl", nil)

	// files in components directory
	createFileFromTemplate(filepath.Join(frontendComponentsPath, "IDForm.js"), "templates/frontend/src/components/IDForm.js.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendComponentsPath, "LoadSubmit.js"), "templates/frontend/src/components/LoadSubmit.js.tmpl", nil)
	createFileFromTemplate(filepath.Join(frontendComponentsPath, "DefaultPage.js"), "templates/frontend/src/components/DefaultPage.js.tmpl", nil)

	// file in public directory
	createFileFromTemplate(filepath.Join(frontendPublicPath, "index.html"), "templates/frontend/public/index.html.tmpl", conf)

	// create App.js
	createFileFromTemplate(filepath.Join(frontendSrcPath, "App.js"), "templates/frontend/src/App.js.tmpl", schemas)

	// create SideBar.js
	createFileFromTemplate(filepath.Join(frontendComponentsPath, "SideBar.js"), "templates/frontend/src/components/SideBar.js.tmpl", schemas)

	// create a schema component and a schema form components for each schema inside schemasWithMethods
	for i := range schemas.List {
		tmpSchema := schemas.List[i]

		// create <schema name>.js components
		createFileFromTemplate(filepath.Join(frontendSchemasPath, tmpSchema.ComponentName+".js"), "templates/frontend/src/components/schemas/Schema.js.tmpl", tmpSchema)

		// create <schema name>Form.js components
		createFileFromTemplate(filepath.Join(frontendSchemaFormsPath, tmpSchema.ComponentName+"Form.js"), "templates/frontend/src/components/schemas/schemaforms/SchemaForm.js.tmpl", tmpSchema)
	}

	log.Info().Msg("Created Frontend successfully.")

	// npm build and moving the output folder +  delete frontend folder
	publishFrontend(frontendPath, filepath.Join(conf.OutputPath, "public"))
}

func publishFrontend(sourcePath string, destinationPath string) {
	// run npm install
	NPMInstall(sourcePath)

	// run npm build
	NPMBuild(sourcePath)

	// move build folder to destinationPath
	fs.MoveDir(filepath.Join(sourcePath, "build"), filepath.Join(destinationPath, "build"))

	// delete frontend folder at source path
	fs.DeleteFolderRecursively(sourcePath)
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

			// add names
			schema.Name = strings.ReplaceAll(strings.ToLower(tmpSchemaName), " ", "")
			schema.H1Name = strings.Title(tmpSchemaName)
			schema.ComponentName = strings.ReplaceAll(schema.H1Name, " ", "")

			// add properties
			schema.Properties = make([]PropertyConf, 0)
			tmpSchemaPropertyNames := reflect.ValueOf(spec.Components.Schemas[tmpSchemaName].Value.Properties).MapKeys()
			for j := range tmpSchemaPropertyNames {
				tmpSchemaPropertyName := tmpSchemaPropertyNames[j].Interface().(string)
				var tmpPropertyConf PropertyConf
				tmpPropertyConf.Name = tmpSchemaPropertyName
				tmpPropertyConf.LabelName = strings.Title(tmpSchemaPropertyName)
				tmpPropertyConf.Type = spec.Components.Schemas[tmpSchemaName].Value.Properties[tmpSchemaPropertyName].Value.Type
				schema.Properties = append(schema.Properties, tmpPropertyConf)
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
	path := filepath.Join(conf.OutputPath, "public", "doc")
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
