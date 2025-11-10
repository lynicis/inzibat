package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReader(t *testing.T) {
	t.Run("json reader", func(t *testing.T) {
		jsonReader, err := NewReaderStrategy(".json")

		assert.NoError(t, err)
		assert.Implements(t, (*ReaderStrategy)(nil), jsonReader)
	})

	t.Run("yaml reader", func(t *testing.T) {
		t.Run("with .yaml extension", func(t *testing.T) {
			yamlReader, err := NewReaderStrategy(".yaml")

			assert.NoError(t, err)
			assert.Implements(t, (*ReaderStrategy)(nil), yamlReader)
		})

		t.Run("with .yml extension", func(t *testing.T) {
			yamlReader, err := NewReaderStrategy(".yml")

			assert.NoError(t, err)
			assert.Implements(t, (*ReaderStrategy)(nil), yamlReader)
		})
	})

	t.Run("toml reader", func(t *testing.T) {
		tomlReader, err := NewReaderStrategy(".toml")

		assert.NoError(t, err)
		assert.Implements(t, (*ReaderStrategy)(nil), tomlReader)
	})

	t.Run("unknown file extension", func(t *testing.T) {
		unknownReader, err := NewReaderStrategy(".unknown")

		assert.Error(t, err)
		assert.Nil(t, unknownReader)
	})
}

func TestReader_ReadConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		cwd, err := getCwd()
		require.NoError(t, err)

		t.Run("json reader", func(t *testing.T) {
			jsonReader := &JsonReader{
				KoanfInstance: koanf.New("."),
			}

			cfg, err := jsonReader.Read(filepath.Join(cwd, "examples", DefaultConfigFileName))

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		})

		t.Run("yaml reader", func(t *testing.T) {
			yamlReader := &YamlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := yamlReader.Read(filepath.Join(cwd, "examples", "inzibat.yaml"))

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		})

		t.Run("toml reader", func(t *testing.T) {
			tomlReader := &TomlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := tomlReader.Read(filepath.Join(cwd, "examples", "inzibat.toml"))

			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		})
	})

	t.Run("config file not found", func(t *testing.T) {
		t.Run("json reader", func(t *testing.T) {
			jsonReader := &JsonReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.Read("")

			assert.Error(t, err, ErrorReadFile)
			assert.Nil(t, cfg)
		})

		t.Run("yaml reader", func(t *testing.T) {
			jsonReader := &YamlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.Read("")

			assert.Error(t, err, ErrorReadFile)
			assert.Nil(t, cfg)
		})

		t.Run("toml reader", func(t *testing.T) {
			jsonReader := &TomlReader{
				KoanfInstance: koanf.New("."),
			}
			cfg, err := jsonReader.Read("")

			assert.Error(t, err, ErrorReadFile)
			assert.Nil(t, cfg)
		})
	})
}

func getCwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, "../"), nil
}
