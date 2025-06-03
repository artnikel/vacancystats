// Package config provides configuration loading from YAML files
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// DBFolderConfig holds configuration for the database folder path
type DBFolderConfig struct {
	DBPath string `yaml:"dbpath"`
}

// LoggingConfig holds configuration for the log file path
type LoggingConfig struct {
	Path string `yaml:"path"`
}

// ResourceConfig holds a list of available job resource names
type ResourceConfig struct {
	ResourceList []string `yaml:"resource_list"`
}

// StatusConfig holds a list of possible job application statuses
type StatusConfig struct {
	StatusList []string `yaml:"status_list"`
}

// Config aggregates all service configurations
type Config struct {
	DBFolder DBFolderConfig `yaml:"dbfolder"`
	Logging  LoggingConfig  `yaml:"logging"`
	Resource ResourceConfig `yaml:"resource"`
	Status   StatusConfig   `yaml:"status"`
}

// LoadConfig loads the configuration from the given YAML file path
func LoadConfig(path string) (*Config, error) {
	// #nosec G304 -- config path is trusted and not user-controlled
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
