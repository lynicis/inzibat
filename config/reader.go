package config

import (
	"errors"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Reader interface {
	ReadConfig(filename string) (*Cfg, error)
}

func NewReader(fileExtension string) (Reader, error) {
	koanfInstance := koanf.New(".")
	readerMapper := map[string]Reader{
		".json": &JsonReader{
			KoanfInstance: koanfInstance,
		},
		".yaml": &YamlReader{
			KoanfInstance: koanfInstance,
		},
		".yml": &YamlReader{
			KoanfInstance: koanfInstance,
		},
		".toml": &TomlReader{
			koanfInstance,
		},
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

func (jsonReader *JsonReader) ReadConfig(filename string) (*Cfg, error) {
	err := jsonReader.KoanfInstance.Load(
		file.Provider(filename),
		json.Parser(),
	)
	if err != nil {
		return nil, errors.New(ErrorReadFile)
	}

	var config Cfg
	err = jsonReader.KoanfInstance.Unmarshal(filename, &config)
	if err != nil {
		return nil, errors.New(ErrorUnmarshalling)
	}

	return &config, nil
}

type YamlReader struct {
	KoanfInstance *koanf.Koanf
}

func (yamlReader *YamlReader) ReadConfig(filename string) (*Cfg, error) {
	err := yamlReader.KoanfInstance.Load(
		file.Provider(filename),
		yaml.Parser(),
	)
	if err != nil {
		return nil, errors.New(ErrorReadFile)
	}

	var config Cfg
	err = yamlReader.KoanfInstance.Unmarshal(filename, &config)
	if err != nil {
		return nil, errors.New(ErrorUnmarshalling)
	}

	return &config, nil
}

type TomlReader struct {
	KoanfInstance *koanf.Koanf
}

func (tomlReader *TomlReader) ReadConfig(filename string) (*Cfg, error) {
	err := tomlReader.KoanfInstance.Load(
		file.Provider(filename),
		toml.Parser(),
	)
	if err != nil {
		return nil, errors.New(ErrorReadFile)
	}

	var config Cfg
	err = tomlReader.KoanfInstance.Unmarshal(filename, &config)
	if err != nil {
		return nil, errors.New(ErrorUnmarshalling)
	}

	return &config, nil
}
