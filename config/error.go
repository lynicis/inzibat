package config

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	ErrorReadFile      = errors.New("error occurred while reading config file")
	ErrorUnmarshalling = errors.New("error occurred while unmarshalling config file")
	ErrorGetSendBody   = errors.New("send body with get http method")
)

// TODO:
func PrettifyValidationError(rawErr validator.InvalidValidationError) error {
	// TODO:
	return nil
}
