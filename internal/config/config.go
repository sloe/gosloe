package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	AppConfig AppConfig `yaml:"app"`
}

func NewConfig() Config {
	return Config{}
}

func NewConfigFromYaml(yamlFilename string) (*Config, error) {
	var config = NewConfig()
	f, err := os.Open(yamlFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)

	return &config, err
}
