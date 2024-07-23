package config

import "github.com/go-playground/validator/v10"

const (
	ErrorReadFile      = "error occurred while reading config file"
	ErrorUnmarshalling = "error occurred while unmarshalling config file"
	ErrorGetSendBody   = "send body with get http method"
)

func PrettifyValidationError(rawErr validator.InvalidValidationError) error {
	// TODO:
	return nil
}
