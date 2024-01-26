package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	configDir  string
	configFile string
)

type Task struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type Config struct {
	Tasks []Task `json:"tasks"`
}

func InitConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir = filepath.Join(homeDir, ".config", "checkmark")
	configFile = filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configDir); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.MkdirAll(configDir, 0751); err != nil {
			return nil, err
		}
	}

	config := Config{}
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &config, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func StoreConfig(config *Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0666)
}
