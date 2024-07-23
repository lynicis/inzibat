package config

import (
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	t.Run("json reader", func(t *testing.T) {
		jsonReader, err := NewReader(".json")

		assert.NoError(t, err)
		assert.Implements(t, (*Reader)(nil), jsonReader)
	})

	t.Run("yaml reader", func(t *testing.T) {
		t.Run("with .yaml extension", func(t *testing.T) {
			yamlReader, err := NewReader(".yaml")

			assert.NoError(t, err)
			assert.Implements(t, (*Reader)(nil), yamlReader)
		})

		t.Run("with .yml extension", func(t *testing.T) {
			yamlReader, err := NewReader(".yml")

			assert.NoError(t, err)
			assert.Implements(t, (*Reader)(nil), yamlReader)
		})
	})

	t.Run("toml reader", func(t *testing.T) {
		tomlReader, err := NewReader(".toml")

		assert.NoError(t, err)
		assert.Implements(t, (*Reader)(nil), tomlReader)
	})

	t.Run("unknown file extension", func(t *testing.T) {
		unknownReader, err := NewReader(".unknown")

		assert.Error(t, err)
		assert.Nil(t, unknownReader)
	})
}

func TestReader_ReadConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		t.Run("json reader", func(t *testing.T) {
			jsonReader := &JsonReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.ReadConfig("../examples/inzibat.config.json")

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		})

		t.Run("yaml reader", func(t *testing.T) {
			jsonReader := &YamlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.ReadConfig("../examples/inzibat.config.yaml")

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		})

		t.Run("toml reader", func(t *testing.T) {
			jsonReader := &TomlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.ReadConfig("../examples/inzibat.config.toml")

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		})
	})

	t.Run("config file not found", func(t *testing.T) {
		t.Run("json reader", func(t *testing.T) {
			jsonReader := &JsonReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.ReadConfig("")

			assert.Errorf(t, err, ErrorReadFile)
			assert.Nil(t, cfg)
		})

		t.Run("yaml reader", func(t *testing.T) {
			jsonReader := &YamlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.ReadConfig("")

			assert.Errorf(t, err, ErrorReadFile)
			assert.Nil(t, cfg)
		})

		t.Run("toml reader", func(t *testing.T) {
			jsonReader := &TomlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.ReadConfig("")

			assert.Errorf(t, err, ErrorReadFile)
			assert.Nil(t, cfg)
		})
	})
}
