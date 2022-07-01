package generator

type Flags struct {
	UseDatabase bool
	UseLogger   bool
	UseHTTP2    bool
}

type AuthConfig struct {
	UseAuth            bool
	ApiKeyHeaderName   string
	ApiKeySecurityName string
}
type GeneratorConfig struct {
	OpenAPIPath  string
	OutputPath   string
	ModuleName   string
	DatabaseName string
	AuthConfig
	Flags
}

type ProjectConfig struct {
	Name string
	Path string
}

type ServerConfig struct {
	Port        int16
	ModuleName  string
	OpenAPIName string
	Flags
}

type ResponseConfig struct {
	StatusCode string
	Desciption string
}

type OperationConfig struct {
	Method      string
	Summary     string
	OperationID string
	UseAuth     bool
	Responses   []ResponseConfig
}

type PathConfig struct {
	Path       string
	Operations []OperationConfig
}

type HandlerConfig struct {
	Paths         []PathConfig
	UseGlobalAuth bool
}

// struct for all schemas that have to be in the frontend
type Schemas struct {
	List       []SchemaConf
	IsNotEmpty bool
}

// struct for the specific schema in Schemas
type SchemaConf struct {
	Name          string            // all in lower case and without spaces
	H1Name        string            // correct grammar, spaces allowed
	ComponentName string            // first letter upper case, no spaces
	Properties    map[string]string // key has first letter in upper case
	Methods       []MethodConf
}

// struct for each method a schema has
type MethodConf struct {
	Type               string
	Endpoint           string
	PathParams         map[string]string
	BodySchemaRequired bool
}
