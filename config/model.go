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
	ServerPort  string
	Routes      []Route
	Concurrency Concurrency
}

type Route struct {
	Method    string
	Path      string
	RequestTo RequestTo
}

type Concurrency struct {
	RouteCreatorLimit int
}

type RequestTo struct {
	Method                 string
	Headers                map[string]string `mapstructure:"headers"`
	Body                   map[string]string `mapstructure:"body"`
	Host                   string
	Path                   string
	PassWithRequestBody    bool
	PassWithRequestHeaders bool
}
