package main

import (
	"errors"

	"github.com/spf13/viper"
)

func ReadConfig() (*Config, error) {
	var err error

	desired := viper.New()
	desired.SetConfigName("desired")
	desired.SetConfigType("json")
	desired.AddConfigPath(".")
	err = desired.ReadInConfig()
	if err != nil {
		return nil, errors.New("error occurred while reading config file")
	}

	var desiredConfig Config
	err = desired.Unmarshal(&desiredConfig)
	if err != nil {
		return nil, errors.New("error occurred while unmarshalling config file")
	}

	return &desiredConfig, nil
}
