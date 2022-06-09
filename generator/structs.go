package generator

import (
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var IMPORT_UUID bool

type TypeConfig struct {
	SchemaDefs  map[string][]TypeDefinition
	ImportDefs  []ImportDefinition
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

func GenerateTypes(spec *openapi3.T, pConf ProjectConfig) {
	schemaDefs := generateStructDefs(&spec.Components.Schemas)
	importDefs := generateImports()
	var conf TypeConfig
	conf.SchemaDefs = schemaDefs
	conf.ImportDefs = importDefs
	conf.ProjectName = pConf.Name

	fileName := "structs.go"
	filePath := filepath.Join(pConf.Path, fileName)
	templateFile := "templates/structs.go.tmpl"
	CreateFileFromTemplate(filePath, templateFile, conf)
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
		var t string
		switch property.Value.Format {
		case "float":
			t = "float32"
		case "int32":
			t = "int32"
		case "uuid":
			IMPORT_UUID = true
			t = "uuid.UUID"
		default:
			t = property.Value.Type
		}
		println(property.Ref)
		if property.Value.Type == "array" {
			if property.Value.Items.Value.Type == "object" {
				t = "[]" + property.Value.Items.Ref
			}
		}
		propertyDef := TypeDefinition{
			name,
			t,
		}
		typeDefs = append(typeDefs, propertyDef)
	}
	return typeDefs
}

func toGoType(sRef *openapi3.SchemaRef) (goType string) {

	// we know the object is defined in the schema
	if sRef.Value.Type == "object" && sRef.Ref != "" {
		goType = strings.Split(sRef.Ref, "/")[0]
	} else if sRef.Value.Type == "array" {
	}
	return goType
}

func generateImports() []ImportDefinition {
	var importDefs []ImportDefinition
	if IMPORT_UUID {
		importDefs = append(importDefs, ImportDefinition{"", "\"github.com/google/uuid\""})
	}

	return importDefs
}
