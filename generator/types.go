package generator

type PortConfig struct {
	Port int16
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
