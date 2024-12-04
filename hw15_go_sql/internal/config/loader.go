package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const defaultConfigFilePath = "configs/config.yaml"

func MustLoadConfig() *Config {
	file, err := os.Open(defaultConfigFilePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	config := &Config{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		panic(err)
	}

	return config
}
