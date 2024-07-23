package config

const (
	EnvironmentVariableConfigFileName = "CONFIG_FN"
	DefaultConfigFileName             = "inzibat.config.json"
)

type Cfg struct {
	ServerPort       int
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
	Headers                map[string][]string
	Body                   map[string]interface{}
	Host                   string
	Path                   string
	PassWithRequestBody    bool
	PassWithRequestHeaders bool
	InErrorReturn500       bool
}

type Mock struct {
	Headers    map[string]string
	Body       map[string]interface{}
	BodyString string
	StatusCode int
}
