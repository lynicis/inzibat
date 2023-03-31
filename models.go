package main

type Config struct {
	ServerPort string   `json:"serverPort"`
	Routes     []Routes `json:"routes"`
}

type Routes struct {
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	RequestTo RequestTo `json:"requestTo"`
}

type RequestTo struct {
	Host string `json:"host"`
	Path string `json:"path"`
}
