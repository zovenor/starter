package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ExecApps []ExecAppConfig `json:"exec_apps"`
}

func New() *Config {
	return new(Config)
}

type ExecAppConfig struct {
	Name     string   `json:"name"`
	Cmds     []string `json:"cmds"`
	StopCmds []string `json:"stop_cmds"`
	Disabled bool     `json:"disabled"`
}

func (config *Config) ImportFromJsonFile(filepath string) error {
	jsonBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, config)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) SaveToJsonFile(filepath string) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
