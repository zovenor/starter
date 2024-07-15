package main

import (
	"github.com/zovenor/logging/prettyPrints"
	"github.com/zovenor/starter/internal/app"
	"github.com/zovenor/starter/internal/flags"
)

func main() {
	prettyPrints.ClearTerminal()
	fl := flags.New()
	err := fl.Parse()
	if err != nil {
		panic(err)
	}

	a, err := app.New(fl.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	if err := a.Run(); err != nil {
		panic(err)
	}
}
