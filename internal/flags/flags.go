package flags

import (
	"flag"
	"os"
	"path"
)

type Flags struct {
	ConfigFilePath string
	HomeDir        string
}

func New() *Flags {
	return new(Flags)
}

func (flags *Flags) validate() error {
	var err error
	if flags.HomeDir == "" {
		flags.HomeDir, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	}
	if flags.ConfigFilePath == "" {
		flags.ConfigFilePath = path.Join(flags.HomeDir, ".starter/config.yaml")
	}
	return nil
}

func (flags *Flags) Parse() error {
	flag.StringVar(&flags.ConfigFilePath, "cfp", "", "config file path")
	flag.StringVar(&flags.HomeDir, "hd", "", "home directory")
	flag.Parse()
	err := flags.validate()
	if err != nil {
		return err
	}
	return nil
}
