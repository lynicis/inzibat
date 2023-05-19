package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

const (
	EnvironmentVariableConfigFileName = "CONFIG_FN"
	DefaultConfigFileName             = "desired"
)

func ReadConfig(filepath, filename string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("json")
	v.AddConfigPath(filepath)

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}

		return nil, errors.New("error occurred while reading config file")
	}

	var desiredConfig Config
	err = v.Unmarshal(&desiredConfig)
	if err != nil {
		return nil, errors.New("error occurred while unmarshalling config file")
	}

	for index, route := range desiredConfig.Routes {
		if route.RequestTo.Method == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if route.Method == http.MethodGet {
			if route.RequestTo.Body != nil {
				return nil, errors.New("send body with get http method")
			}
		}

		desiredConfig.Routes[index].Method = route.Method
		desiredConfig.Routes[index].RequestTo.Method = route.RequestTo.Method
	}

	return &desiredConfig, nil
}

func (c *Config) Print() {
	fmt.Println("🫡 INZIBAT 🪖")
	fmt.Println(
		fmt.Sprintf(
			"Open Routes: %d", len(c.Routes),
		),
	)
	fmt.Println(
		fmt.Sprintf(
			"Server Port: %s", c.ServerPort,
		),
	)
}
