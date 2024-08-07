package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ExecApps []ExecAppConfig `yaml:"exec_apps"`
	Vars     []Var           `yaml:"vars"`
}

type Var struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required"`
	Hiden    bool   `yaml:"hiden"`
}

func New() *Config {
	return new(Config)
}

type ExecAppConfig struct {
	Name     string   `yaml:"name"`
	Cmds     []string `yaml:"cmds"`
	StopCmds []string `yaml:"stop_cmds"`
	Disabled bool     `yaml:"disabled"`
	CheckCmd string   `yaml:"check_cmd"`
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
