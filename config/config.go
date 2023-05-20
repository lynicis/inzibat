package config

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func ReadConfig(filepath, filename string) (*Config, error) {
	viperInstance := viper.New()
	viperInstance.SetConfigName(filename)
	viperInstance.SetConfigType(FileTypeJson)
	viperInstance.AddConfigPath(filepath)

	err := viperInstance.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New(ErrorFileNotFound)
		}

		return nil, errors.New(ErrorReadFile)
	}

	var desiredConfig *Config
	err = viperInstance.Unmarshal(&desiredConfig)
	if err != nil {
		return nil, errors.New(ErrorUnmarshalling)
	}

	for index, route := range desiredConfig.Routes {
		if route.RequestTo.Method == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if route.Method == http.MethodGet {
			if route.RequestTo.Body != nil {
				return nil, errors.New(ErrorGetSendBody)
			}
		}

		desiredConfig.Routes[index].Method = route.Method
		desiredConfig.Routes[index].RequestTo.Method = route.RequestTo.Method
	}

	return desiredConfig, nil
}

func (c *Config) Print() {
	fmt.Println("🫡 INZIBAT 🪖")
	fmt.Printf(
		"Open Routes: %d\n", len(c.Routes),
	)
	fmt.Printf(
		"Server Port: %s", c.ServerPort,
	)
}
