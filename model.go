package main

type HttpHeader map[string]string
type HttpBody map[string]any

type HttpResponse struct {
	Status int
	Body   []byte
}

type Config struct {
	ServerPort string  `json:"serverPort"`
	Routes     []Route `json:"routes"`
}

type Route struct {
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	RequestTo RequestTo `json:"requestTo"`
}

type RequestTo struct {
	Method string     `json:"method,omitempty"`
	Header HttpHeader `json:"header,omitempty"`
	Body   HttpBody   `json:"body,omitempty"`
	Host   string     `json:"host"`
	Path   string     `json:"path"`
}
