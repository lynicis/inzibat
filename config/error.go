package config

import (
	"errors"
	"fmt"
)

var (
	ErrorReadFile      = errors.New("error occurred while reading config file")
	ErrorUnmarshalling = errors.New("error occurred while unmarshalling config file")
	ErrorGetSendBody   = errors.New("send body with get http method")
)

func newFailOpeningError(err error) error {
	return fmt.Errorf("failed to open file: %w", err)
}

func newFailReadingError(err error) error {
	return fmt.Errorf("failed to read file: %w", err)
}
