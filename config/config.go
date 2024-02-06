package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kevhlee/checkmark/task"
)

var (
	configDir  string
	configFile string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configDir = filepath.Join(homeDir, ".config", "checkmark")
	configFile = filepath.Join(configDir, "config.json")
}

type Config struct {
	Tasks []task.Task `json:"tasks"`
}

func Init() (Config, error) {
	var cfg Config

	if _, err := os.Stat(configDir); err != nil {
		if !os.IsNotExist(err) {
			return cfg, err
		}
		if err := os.MkdirAll(configDir, 0751); err != nil {
			return cfg, err
		}
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func Save(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0666)
}
