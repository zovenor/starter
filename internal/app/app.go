package app

import "github.com/zovenor/starter/internal/config"

type App struct {
	Name           string
	Version        string
	Config         *config.Config
	configFilePath string
}

func New(configFilePath, name, version string) (*App, error) {
	a := new(App)
	a.Name = name
	a.Version = version
	a.configFilePath = configFilePath
	a.Config = config.New()
	err := a.Config.ImportFromYamlFile(configFilePath)
	if err != nil {
		return nil, err
	}
	return a, nil
}
