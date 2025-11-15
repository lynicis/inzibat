package form_builder

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ValidateStatusCode(codeStr string) error {
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return fmt.Errorf("invalid status code")
	}

	if code < http.StatusContinue || code > http.StatusNetworkAuthenticationRequired {
		return fmt.Errorf("status code must be between 100 and 511")
	}

	return nil
}

func ValidatePath(path string) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("route path must start with '/'")
	}

	return nil
}

func ValidateHost(host string) error {
	if _, err := url.Parse(host); err != nil {
		return fmt.Errorf("invalid hostname")
	}

	return nil
}

func ValidateNonEmpty(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	return nil
}
