package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFailOpeningError(t *testing.T) {

	baseErr := errors.New("original error")

	err := newFailOpeningError(baseErr)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")
	assert.ErrorIs(t, err, baseErr)
}

func TestNewFailReadingError(t *testing.T) {
	baseErr := errors.New("original read error")

	err := newFailReadingError(baseErr)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
	assert.ErrorIs(t, err, baseErr)
}
