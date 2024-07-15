package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/starter/internal/models/major"
)

func (a *App) Run() error {
	mm := major.New()
	mm.Config = a.Config
	mm.Name = a.Name
	mm.Version = a.Version
	mm.ConfigFilePath = a.configFilePath
	p := tea.NewProgram(mm)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
