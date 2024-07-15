package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ExecApps []ExecAppConfig `yaml:"exec_apps"`
}

func New() *Config {
	return new(Config)
}

type ExecAppConfig struct {
	Name     string   `yaml:"name"`
	Cmds     []string `yaml:"cmds"`
	StopCmds []string `yaml:"stop_cmds"`
	Disabled bool     `yaml:"disabled"`
}

func (config *Config) ImportFromYamlFile(filepath string) error {
	filebytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(filebytes, config)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) SaveToYamlFile(filepath string) error {
	filebytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, filebytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
