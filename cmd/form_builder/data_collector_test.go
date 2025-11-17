package form_builder

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"inzibat/config"
)

func TestCollectHeaders_PublicFunction(t *testing.T) {
	t.Run("happy path - CollectHeaders function exists and returns correct type", func(t *testing.T) {
		var _ func() (http.Header, error) = CollectHeaders

		assert.NotNil(t, CollectHeaders)
	})
}

func TestCollectBody_PublicFunction(t *testing.T) {
	t.Run("happy path - CollectBody function exists and returns correct type", func(t *testing.T) {
		var _ func() (config.HttpBody, error) = CollectBody

		assert.NotNil(t, CollectBody)
	})
}

func TestCollectBodyString_PublicFunction(t *testing.T) {
	t.Run("happy path - CollectBodyString function exists and returns correct type", func(t *testing.T) {
		var _ func() (string, error) = CollectBodyString

		assert.NotNil(t, CollectBodyString)
	})
}
