package generator

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
)

var IMPORT_UUID bool
var IMPORT_TIME bool

type ModelCOnfig struct {
	Imports     ImportsConfig
	SchemaDefs  map[string][]TypeDefinition
	ProjectName string
}

type TypeDefinition struct {
	Name        string
	Type        string
	MarshalName string
	// only if Type is struct
	NestedTypes []TypeDefinition
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
	var conf ModelCOnfig
	conf.Imports = imports
	conf.SchemaDefs = schemaDefs
	conf.ProjectName = pConf.Name

	fileName := "model.go"
	filePath := filepath.Join(pConf.Path, Pkg, ModelPkg, fileName)
	templateFiles := []string{"templates/model.go.tmpl", "templates/imports.go.tmpl", "templates/structs.go.tmpl"}
	createFileFromTemplates(filePath, templateFiles, conf)
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
	i := 0
	for name, property := range *properties {
		goType, nested := toGoType(property)
		var nestedGoTypes []TypeDefinition
		if nested {
			nestedGoTypes = generateTypeDefs(&property.Value.Properties)
		}

		// first letter to lower case
		marshalName := []rune(name)
		marshalName[0] = unicode.ToLower(marshalName[0])
		propertyDef := TypeDefinition{
			name,
			goType,
			string(marshalName),
			nestedGoTypes,
		}
		typeDefs[i], i = propertyDef, i+1
	}

	return typeDefs
}

// schema type to generated go type
func toGoType(sRef *openapi3.SchemaRef) (goType string, nested bool) {

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
		items, _ := toGoType(sRef.Value.Items)
		goType = "[]" + items
	case "object":
		if sRef.Ref != "" {
			// checks if object type is defined by reference elsewhere in the schema
			splitRef := strings.Split(sRef.Ref, "/")
			goType = splitRef[len(splitRef)-1]
		} else {
			goType = "struct"
			nested = true
		}
	default:
		goType = sRef.Value.Type
	}
	return goType, nested
}

func generateImports() ImportsConfig {
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

	return conf
}
