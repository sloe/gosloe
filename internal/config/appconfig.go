package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	TreeRoot string `yaml:"treeroot"`
}

func NewAppConfig() AppConfig {
	return AppConfig{}
}

func NewAppConfigFromYaml(yamlFilename string) (*AppConfig, error) {
	var appConfig = NewAppConfig()
	f, err := os.Open(yamlFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(appConfig)
	return &appConfig, err
}
