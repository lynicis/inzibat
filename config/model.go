package config

import (
	"github.com/Lynicis/inzibat/client"
)

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
	Method string            `json:"method,omitempty"`
	Header client.HttpHeader `json:"header,omitempty"`
	Body   []byte            `json:"body,omitempty"`
	Host   string            `json:"host"`
	Path   string            `json:"path"`
}
