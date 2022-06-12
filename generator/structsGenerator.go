package generator

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
)

var IMPORT_UUID bool
var IMPORT_TIME bool

type TypeConfig struct {
	Imports     string
	SchemaDefs  map[string][]TypeDefinition
	ProjectName string
}

type TypeDefinition struct {
	Name string
	Type string
}

type ImportDefinition struct {
	Name string
	URL  string
}

type ImportsConfig struct {
	ImportDefs []ImportDefinition
}

func GenerateTypes(spec *openapi3.T, pConf ProjectConfig) {
	schemaDefs := generateStructDefs(&spec.Components.Schemas)
	imports := generateImports()
	var conf TypeConfig
	conf.Imports = imports
	conf.SchemaDefs = schemaDefs
	conf.ProjectName = pConf.Name

	fileName := "structs.go"
	filePath := filepath.Join(pConf.Path, Pkg, ModelPkg, fileName)
	templateFile := "templates/structs.go.tmpl"
	createFileFromTemplate(filePath, templateFile, conf)
}

func generateStructDefs(schemas *openapi3.Schemas) map[string][]TypeDefinition {
	schemaDefs := make(map[string][]TypeDefinition, len(*schemas))

	for schemaName, ref := range *schemas {
		schemaDefs[schemaName] = generateTypeDefs(&ref.Value.Properties)
	}
	return schemaDefs
}

func generateTypeDefs(properties *openapi3.Schemas) []TypeDefinition {
	typeDefs := make([]TypeDefinition, len(*properties))

	for name, property := range *properties {
		propertyDef := TypeDefinition{
			name,
			toGoType(property),
		}
		typeDefs = append(typeDefs, propertyDef)
	}
	return typeDefs
}

// schema type to generated go type
func toGoType(sRef *openapi3.SchemaRef) (goType string) {

	switch sRef.Value.Type {
	case "number":
		switch sRef.Value.Format {
		case "float":
			goType = "float32"
		case "double":
			goType = "float64"
		default:
			goType = "float"
		}
	case "integer":
		goType = sRef.Value.Format
	case "string":
		switch sRef.Value.Format {
		case "binary":
			goType = "[]byte"
		case "date":
			IMPORT_TIME = true
			goType = "time.Time"
		case "uuid":
			IMPORT_UUID = true
			goType = "uuid.UUID"
		default:
			goType = "string"
		}
	case "array":
		goType = "[]" + toGoType(sRef.Value.Items)
	case "object":
		if sRef.Ref != "" {
			// checks if object type is defined by reference elsewhere in the schema
			splitRef := strings.Split(sRef.Ref, "/")
			goType = splitRef[len(splitRef)-1]

		} else {
			// TODO nested structs
		}
	default:
		goType = sRef.Value.Type
	}
	return goType
}

func generateImports() string {
	var importDefs []ImportDefinition
	if IMPORT_UUID {
		importDefs = append(importDefs, ImportDefinition{"", "\"github.com/google/uuid\""})
	}
	if IMPORT_TIME {
		importDefs = append(importDefs, ImportDefinition{"time", ""})
	}

	conf := ImportsConfig{
		importDefs,
	}

	templateFile := "templates/imports.go.tmpl"
	buf := &bytes.Buffer{}

	tmpl := template.Must(template.ParseFiles(templateFile))
	if tmplErr := tmpl.Execute(buf, conf); tmplErr != nil {
		log.Fatal().Err(tmplErr).Msg("Failed executing imports template.")
	}

	return buf.String()
}
