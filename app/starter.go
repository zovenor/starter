package main

import (
	"os"

	"github.com/zovenor/logging/prettyPrints"
	"github.com/zovenor/logging/v2"
	"github.com/zovenor/starter/internal/app"
	"github.com/zovenor/starter/internal/flags"
)

var version string

const name = "Starter"

func main() {
	prettyPrints.ClearTerminal()
	fl := flags.New()
	err := fl.Parse()
	if err != nil {
		logging.Fatal(err)
		os.Exit(1)
	}

	a, err := app.New(fl.ConfigFilePath, name, version)
	if err != nil {
		logging.Fatal(err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		logging.Fatal(err)
		os.Exit(1)

	}
}
