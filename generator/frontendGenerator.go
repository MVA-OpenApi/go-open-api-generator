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

	// TODO add methods to these schemas

	fmt.Println(schemas)

	log.Info().Msg("Created Frontend.")
}

/* func addMethods(spec *openapi3.T, schemas Schemas) {

} */

func createSchemas(spec *openapi3.T) (schemas Schemas) {
	schemas.List = make([]SchemaConf, 0)

	schemaStrings := toString(reflect.ValueOf(spec.Components.Schemas).MapKeys())
	//pathStrings := toString(reflect.ValueOf(spec.Paths).MapKeys())

	for i := range schemaStrings {
		tmpSchemaName := schemaStrings[i]

		// check if schema has x-label == "form" -> if yes add schema to list
		schemaInformation, _ := spec.Components.Schemas[tmpSchemaName].Value.MarshalJSON()
		if strings.Contains(string(schemaInformation[:]), "\"x-label\":\"form\"") {
			var schema SchemaConf
			schema.Properties = make(map[string]string)

			// add name
			schema.Name = tmpSchemaName

			// add properties
			tmpSchemaPropertyNames := reflect.ValueOf(spec.Components.Schemas[tmpSchemaName].Value.Properties).MapKeys()
			for j := range tmpSchemaPropertyNames {
				tmpSchemaPropertyName := tmpSchemaPropertyNames[j].Interface().(string)
				schema.Properties[tmpSchemaPropertyName] = spec.Components.Schemas[tmpSchemaName].Value.Properties[tmpSchemaPropertyName].Value.Type
			}

			schemas.List = append(schemas.List, schema)
		}

		//test, _ := spec.Components.Schemas[schemaStrings[i]].JSONLookup("x-label")
		//fmt.Println(test)
		/* if schemaProperties != nil {
			fmt.Println(schemaStrings[i] + "has an x-label")
		} */

		//fmt.Println(string(schemaProperties)[:])

		// add properties spec.Components.Schemas[schemaStrings[i]].Value.Properties
		/* schemaProperties := reflect.ValueOf(spec.Components.Schemas[schemaStrings[i]].Value.Properties).MapKeys()
		for i := range schemaProperties {

		} */

	}

	//fmt.Println(schemaStrings)
	//fmt.Println(pathStrings)

	/* jsonSchema, _ := spec.Components.Schemas["Store"].MarshalJSON()
	fmt.Println(string(jsonSchema)[:]) */

	/* jsonComponents, _ := spec.Components.MarshalJSON()
	fmt.Println("Components:")
	fmt.Println(string(jsonComponents)[:])

	jsonPaths, _ := spec.Paths["/store"].MarshalJSON()
	fmt.Println("Paths:")
	fmt.Println(string(jsonPaths)[:]) */

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
