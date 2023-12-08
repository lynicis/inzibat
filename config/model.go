package config

const (
	ErrorFileNotFound  = "config file not found"
	ErrorReadFile      = "error occurred while reading config file"
	ErrorUnmarshalling = "error occurred while unmarshalling config file"
	ErrorGetSendBody   = "send body with get http method"
)

const (
	EnvironmentVariableConfigFileName = "CONFIG_FN"
	DefaultConfigFileName             = "inzibat.config.json"
)

type Config struct {
	ServerPort       string
	Routes           []Route
	Concurrency      Concurrency
	HealthCheckRoute bool
}

type Route struct {
	Method    string
	Path      string
	RequestTo RequestTo
	Mock      Mock
}

type Concurrency struct {
	RouteCreatorLimit int
}

type RequestTo struct {
	Method                 string
	Headers                map[string]string
	Body                   map[string]interface{}
	Host                   string
	Path                   string
	PassWithRequestBody    bool
	PassWithRequestHeaders bool
	InErrorReturn500       bool
}

type Mock struct {
	Headers map[string]string
	Body    map[string]interface{}
	Status  int
}
