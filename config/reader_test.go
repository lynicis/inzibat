package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReader(t *testing.T) {
	for _, extension := range []string{".json", ".yaml", ".yml", ".toml"} {
		t.Run(extension, func(t *testing.T) {
			reader, err := NewReaderStrategy(extension)

			assert.NoError(t, err)
			assert.Implements(t, (*ReaderStrategy)(nil), reader)
		})
	}

	t.Run("unknown file extension", func(t *testing.T) {
		unknownReader, err := NewReaderStrategy(".unknown")

		assert.Error(t, err)
		assert.Nil(t, unknownReader)
	})
}

func TestReader_ReadConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		testCases := map[string]string{
			".json": "inzibat.json",
			".yaml": "inzibat.yaml",
			".yml":  "inzibat.yaml",
			".toml": "inzibat.toml",
		}

		for extension, fileName := range testCases {
			t.Run(extension, func(t *testing.T) {
				reader, err := NewReaderStrategy(extension)
				require.NoError(t, err)

				cfg, err := reader.Read(filepath.Join("../", "examples", fileName))

				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			})
		}
	})

	t.Run("config file not found", func(t *testing.T) {
		reader, err := NewReaderStrategy(DefaultConfigExtension)
		require.NoError(t, err)

		cfg, err := reader.Read("")

		assert.ErrorIs(t, err, ErrorReadFile)
		assert.Nil(t, cfg)
	})

	t.Run("config unmarshalling error", func(t *testing.T) {
		testCases := map[string]string{
			".json": "bad_unmarshal.json",
			".yaml": "bad_unmarshal.yaml",
			".toml": "bad_unmarshal.toml",
		}

		for extension, fileName := range testCases {
			t.Run(extension, func(t *testing.T) {
				reader, err := NewReaderStrategy(extension)
				require.NoError(t, err)

				cfg, err := reader.Read(fileName)

				assert.ErrorIs(t, err, ErrorUnmarshalling)
				assert.Nil(t, cfg)
			})
		}
	})
}

func TestReader_ReadDoesNotAccumulateState(t *testing.T) {
	first := filepath.Join(t.TempDir(), "first.json")
	require.NoError(t, os.WriteFile(first, []byte(`{"serverPort":1234}`), 0o600))

	second := filepath.Join(t.TempDir(), "second.json")
	require.NoError(t, os.WriteFile(second, []byte(`{"routes":[]}`), 0o600))

	reader, err := NewReaderStrategy(DefaultConfigExtension)
	require.NoError(t, err)

	firstCfg, err := reader.Read(first)
	require.NoError(t, err)
	assert.Equal(t, 1234, firstCfg.ServerPort)

	secondCfg, err := reader.Read(second)
	require.NoError(t, err)
	assert.Zero(t, secondCfg.ServerPort)
}
