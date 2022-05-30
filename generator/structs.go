package generator

import (
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

type TypeConfig struct {
	SchemaDefs  map[string][]TypeDefinition
	ProjectName string
}

type TypeDefinition struct {
	TypeName string
	Type     string
}

func GenerateTypes(spec *openapi3.T, pConf ProjectConfig) {
	schemaDefs := generateStructDefs(&spec.Components.Schemas)
	var conf TypeConfig
	conf.SchemaDefs = schemaDefs
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
			/* case "uuid":
				t = openapi_types.UUID
			} */
		default:
			t = property.Value.Type
		}
		propertyDef := TypeDefinition{
			name,
			t,
		}
		typeDefs = append(typeDefs, propertyDef)
	}
	return typeDefs
}
