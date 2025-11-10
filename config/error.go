package config

import "errors"

var (
	ErrorReadFile      = errors.New("error occurred while reading config file")
	ErrorUnmarshalling = errors.New("error occurred while unmarshalling config file")
	ErrorGetSendBody   = errors.New("send body with get http method")
)
