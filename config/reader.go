package config

import (
	"errors"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type ReaderStrategy interface {
	Read(filename string) (*Cfg, error)
}

type koanfReader struct {
	parser koanf.Parser
}

func NewReaderStrategy(fileExtension string) (ReaderStrategy, error) {
	parsers := map[string]koanf.Parser{
		DefaultConfigExtension: json.Parser(),
		".yaml":                yaml.Parser(),
		".yml":                 yaml.Parser(),
		".toml":                toml.Parser(),
	}

	parser := parsers[fileExtension]
	if parser == nil {
		return nil, errors.New("ambiguous file extension")
	}

	return &koanfReader{parser: parser}, nil
}

func (reader *koanfReader) Read(filename string) (*Cfg, error) {
	instance := koanf.New(".")

	if err := instance.Load(file.Provider(filename), reader.parser); err != nil {
		return nil, ErrorReadFile
	}

	var config *Cfg
	if err := instance.Unmarshal("", &config); err != nil {
		return nil, ErrorUnmarshalling
	}

	return config, nil
}
