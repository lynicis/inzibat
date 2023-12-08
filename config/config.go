package config

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

func ReadConfig(filename string) (*Config, error) {
	viperInstance := viper.New()
	viperInstance.SetConfigName(filename)

	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, errors.New(ErrorReadFile)
	}

	cleanFilename := path.Clean(filename)
	configFilePath := path.Join(workingDirectory, cleanFilename)
	extensionOfFilePath := path.Ext(configFilePath)
	if extensionOfFilePath == "" {
		configFilePath = path.Clean(configFilePath + ".json")
	}

	viperInstance.SetConfigFile(configFilePath)
	err = viperInstance.ReadInConfig()
	if err != nil {
		var errorConfigFileNotFound *fs.PathError
		isErrorConfigFileNotFound := errors.As(err, &errorConfigFileNotFound)
		if isErrorConfigFileNotFound {
			return nil, errors.New(ErrorFileNotFound)
		}

		return nil, errors.New(ErrorReadFile)
	}

	/*
	 * case-insensitive map keys
	 */
	var desiredConfig Config
	err = viperInstance.Unmarshal(&desiredConfig)
	if err != nil {
		return nil, errors.New(ErrorUnmarshalling)
	}

	for indexOfRoute, route := range desiredConfig.Routes {
		if route.RequestTo.Method == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if route.RequestTo.Method == http.MethodGet {
			if route.RequestTo.Body != nil {
				return nil, errors.New(ErrorGetSendBody)
			}
		}

		desiredConfig.Routes[indexOfRoute].Method = route.Method
		desiredConfig.Routes[indexOfRoute].RequestTo.Method = route.RequestTo.Method
	}

	if desiredConfig.Concurrency.RouteCreatorLimit == 0 {
		desiredConfig.Concurrency.RouteCreatorLimit = runtime.NumCPU()
	}

	return &desiredConfig, nil
}
