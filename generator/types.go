package generator

type GeneratorConfig struct {
	OpenAPIPath string
	OutputPath  string
	ModuleName  string
	UseLogger   bool
	OpenAPIName string
}

type ProjectConfig struct {
	Name string
	Path string
}

type ServerConfig struct {
	Port        int16
	ModuleName  string
	UseLogger   bool
	OpenAPIName string
}

type ResponseConfig struct {
	StatusCode string
	Desciption string
}

type OperationConfig struct {
	Method      string
	Summary     string
	OperationID string
	Responses   []ResponseConfig
}

type PathConfig struct {
	Path       string
	Operations []OperationConfig
}

type HandlerConfig struct {
	Paths []PathConfig
}
