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

func NewReaderStrategy(fileExtension string) (ReaderStrategy, error) {
	koanfInstance := koanf.New(".")
	readerMapper := map[string]ReaderStrategy{
		".json": &JsonReader{KoanfInstance: koanfInstance},
		".yaml": &YamlReader{KoanfInstance: koanfInstance},
		".yml":  &YamlReader{KoanfInstance: koanfInstance},
		".toml": &TomlReader{KoanfInstance: koanfInstance},
	}

	reader := readerMapper[fileExtension]
	if reader == nil {
		return nil, errors.New("ambiguous file extension")
	}

	return reader, nil
}

type JsonReader struct {
	KoanfInstance *koanf.Koanf
}

func (jsonReader *JsonReader) Read(filename string) (*Cfg, error) {
	if err := jsonReader.KoanfInstance.Load(
		file.Provider(filename),
		json.Parser(),
	); err != nil {
		return nil, ErrorReadFile
	}

	var config *Cfg
	if err := jsonReader.KoanfInstance.Unmarshal("", &config); err != nil {
		return nil, ErrorUnmarshalling
	}

	return config, nil
}

type YamlReader struct {
	KoanfInstance *koanf.Koanf
}

func (yamlReader *YamlReader) Read(filename string) (*Cfg, error) {
	if err := yamlReader.KoanfInstance.Load(
		file.Provider(filename),
		yaml.Parser(),
	); err != nil {
		return nil, ErrorReadFile
	}

	var config *Cfg
	if err := yamlReader.KoanfInstance.Unmarshal("", &config); err != nil {
		return nil, ErrorUnmarshalling
	}

	return config, nil
}

type TomlReader struct {
	KoanfInstance *koanf.Koanf
}

func (tomlReader *TomlReader) Read(filename string) (*Cfg, error) {
	if err := tomlReader.KoanfInstance.Load(
		file.Provider(filename),
		toml.Parser(),
	); err != nil {
		return nil, ErrorReadFile
	}

	var config *Cfg
	if err := tomlReader.KoanfInstance.Unmarshal("", &config); err != nil {
		return nil, ErrorUnmarshalling
	}

	return config, nil
}
