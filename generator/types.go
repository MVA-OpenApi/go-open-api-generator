package generator

type Flags struct {
	UseDatabase   bool
	UseLogger     bool
	UseHTTP2      bool
	UseValidation bool
	UseLifecycle  bool
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

type Step struct {
	Method        string
	Endpoint      string
	Payload       string
	Name          string
	StatusCode    string
	RegexPaths    []string
	RealName      string //string that will be used in InitializeScenario
	RegexAndCode  map[string]int
	PathsWithHost []string
	Mapping       map[string]int //it may occur that we use the same function in order to reach different endpoints
}

type Listing struct {
	Steps           []Step
	UniqueEndpoints []string
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
	OpenAPIPath   string
	UseAuth       bool
	UseGlobalAuth bool
	ModuleName    string
	Flags
}

// struct for all schemas that have to be in the frontend
type Schemas struct {
	List       []SchemaConf
	IsNotEmpty bool
}

// struct for the specific schema in Schemas
type SchemaConf struct {
	Name          string // all in lower case and without spaces
	H1Name        string // correct grammar, spaces allowed
	ComponentName string // first letter upper case, no spaces
	Properties    []PropertyConf
	Methods       []MethodConf
}

// struct for each property of a schema
type PropertyConf struct {
	Name      string
	LabelName string
	Type      string
}

// struct for each method a schema has
type MethodConf struct {
	Type               string
	Endpoint           string
	PathParams         map[string]string
	BodySchemaRequired bool
}
