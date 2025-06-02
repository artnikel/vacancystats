package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DBFolderConfig struct {
	DBPath string `yaml:"dbpath"`
}

type LoggingConfig struct {
	Path string `yaml:"path"`
}

type ResourceConfig struct {
	ResourceList []string `yaml:"resource_list"`
}

type StatusConfig struct {
	StatusList []string `yaml:"status_list"`
}

type Config struct {
	DBFolder DBFolderConfig `yaml:"dbfolder"`
	Logging  LoggingConfig  `yaml:"logging"`
	Resource ResourceConfig `yaml:"resource"`
	Status   StatusConfig   `yaml:"status"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
