package major

import (
	"fmt"
	"os/user"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/starter/internal/config"
	"github.com/zovenor/starter/internal/executable"
	"github.com/zovenor/starter/internal/keymap"
	"github.com/zovenor/starter/internal/models/vars"
)

type MajorModel struct {
	Name           string
	Version        string
	Config         *config.Config
	ExecutableApps []*executable.ExecutableApp
	ConfigFilePath string

	cursor int
	help   help.Model
	keys   keymap.KeyMap

	height, width int
	currentTime   time.Time
	username      string

	varsModel *vars.VarsModel
}

func New() *MajorModel {
	mm := new(MajorModel)
	mm.help = help.New()
	mm.keys = keymap.Keys
	mm.height = 30
	mm.width = 10000
	mm.currentTime = time.Now()
	cu, err := user.Current()
	if err != nil {
		panic(err)
	}
	mm.username = cu.Name
	return mm
}

func (mm *MajorModel) sortExecApps() {
	execAppNewList := make([]config.ExecAppConfig, 0, len(mm.Config.ExecApps))
	for _, execAppConfig := range mm.Config.ExecApps {
		if execAppConfig.Disabled {
			continue
		}
		execAppNewList = append(execAppNewList, execAppConfig)
	}
	for _, execAppConfig := range mm.Config.ExecApps {
		if execAppConfig.Disabled {
			execAppNewList = append(execAppNewList, execAppConfig)
		}
	}
	mm.Config.ExecApps = execAppNewList
	mm.setExecutableApps()
}

func (mm *MajorModel) GetExecAppByName(name string) (*executable.ExecutableApp, error) {
	for _, execApp := range mm.ExecutableApps {
		if execApp.Name == name {
			return execApp, nil
		}
	}
	return nil, fmt.Errorf("can not find executable app with name %v", name)
}

func (mm *MajorModel) setExecutableApps() {
	if len(mm.ExecutableApps) == 0 {
		mm.ExecutableApps = make([]*executable.ExecutableApp, 0, len(mm.Config.ExecApps))
		for i, execApp := range mm.Config.ExecApps {
			mm.ExecutableApps = append(mm.ExecutableApps, &executable.ExecutableApp{
				ExecAppConfig: execApp,
				Index:         i,
			})
		}
	} else {
		newExecutableAppsList := make([]*executable.ExecutableApp, 0, len(mm.Config.ExecApps))
		for i, execAppCfg := range mm.Config.ExecApps {
			execApp, err := mm.GetExecAppByName(execAppCfg.Name)
			if err != nil {
				continue
			}
			execApp.ExecAppConfig = execAppCfg
			execApp.Index = i
			newExecutableAppsList = append(newExecutableAppsList, execApp)
		}
		mm.ExecutableApps = newExecutableAppsList
	}
}

type TickMsg time.Time

func (mm *MajorModel) doTick() tea.Cmd {
	return tea.Tick(time.Second, func(ct time.Time) tea.Msg {
		mm.currentTime = ct
		return TickMsg(ct)
	})
}

func (mm *MajorModel) Init() tea.Cmd {
	mm.setExecutableApps()
	mm.varsModel = vars.New(mm.Config.Vars, mm, mm.keys, mm.help, mm.String())
	cmds := make([]tea.Cmd, 0)
	for _, execApp := range mm.ExecutableApps {
		if cmd := mm.CheckExecApp(execApp, 10*time.Second); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	cmds = append(cmds, mm.doTick())
	return tea.Batch(cmds...)
}

func (mm *MajorModel) RevertDisabled(index int) {
	if mm.ExecutableApps[index].Status() == executable.IsNotRunning {
		mm.Config.ExecApps[index].Disabled = !mm.Config.ExecApps[index].Disabled
		mm.sortExecApps()
	}
	mm.Config.SaveToYamlFile(mm.ConfigFilePath)
}

func (mm *MajorModel) Cursor() int {
	return mm.cursor
}

func (mm *MajorModel) NextCursor() {
	if mm.cursor < len(mm.ExecutableApps)-1 {
		mm.cursor++
	}
}

func (mm *MajorModel) LastCursor() {
	if mm.cursor > 0 {
		mm.cursor--
	}
}

func (mm *MajorModel) String() string {
	return fmt.Sprintf("%v %v", mm.Name, mm.Version)
}

func (mm *MajorModel) View() string {
	var s string

	timeView := fmt.Sprintf("%v", mm.currentTime.Format(time.Stamp))
	timeView = fmt.Sprintf("%v  %v", mm.username, timeView)
	s += mm.String()
	s += strings.Repeat(" ", mm.width-len(s)-len(timeView))
	s += timeView
	s += "\n\n"

	separated := false
	for _, execApp := range mm.ExecutableApps {
		if execApp.Disabled && !separated {
			s += "\nDetached:\n"
			separated = true
		}
		s += execApp.Format(execApp.Index == mm.Cursor())
	}

	keysView := mm.help.View(mm.keys)
	keysView += "\n"
	rpt := mm.height - strings.Count(s, "\n") - strings.Count(keysView, "\n") - 1
	if rpt < 0 {
		rpt = 0
	}
	s += strings.Repeat("\n", rpt)
	s += keysView

	return s
}

func (mm *MajorModel) CurrentExecApp() (*executable.ExecutableApp, error) {
	if len(mm.ExecutableApps) == 0 {
		return nil, fmt.Errorf("len of app is zero")
	}
	return mm.ExecutableApps[mm.cursor], nil
}

func (mm *MajorModel) RunExecApp(execApp *executable.ExecutableApp) tea.Cmd {
	go execApp.Run(mm.varsModel.Vars)
	return nil
}

func (mm *MajorModel) CheckExecApp(execApp *executable.ExecutableApp, timeout time.Duration) tea.Cmd {
	go func() {
		for {
			execApp.Check()
			if timeout == 0 {
				return
			}
			time.Sleep(timeout)
		}
	}()
	return checkExecutableApp(execApp)
	return nil
}

func (mm *MajorModel) StopExecApp(execApp *executable.ExecutableApp) tea.Cmd {
	go execApp.Stop()
	return nil
}

func (mm *MajorModel) checkRequiredVar(withCmd tea.Cmd) (tea.Model, tea.Cmd) {
	for _, vwv := range mm.varsModel.Vars {
		if vwv.Required && vwv.Value == "" {
			mm.varsModel.WithCmd = withCmd
			return mm.varsModel, func() tea.Msg {
				return tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'e'},
				}
			}
		}
	}
	return nil, nil
}

func (mm *MajorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, mm.keys.Quit):
			return mm, tea.Quit
		case key.Matches(msg, mm.keys.Up):
			mm.LastCursor()
			return mm, nil
		case key.Matches(msg, mm.keys.Down):
			mm.NextCursor()
			return mm, nil
		case key.Matches(msg, mm.keys.Detach):
			execApp, err := mm.CurrentExecApp()
			if err != nil {
				return mm, nil
			}
			mm.RevertDisabled(execApp.Index)
			return mm, nil
		case key.Matches(msg, mm.keys.Run):
			m, cmd := mm.checkRequiredVar(func() tea.Msg { return msg })
			if m != nil {
				return m, cmd
			}
			currentExecApp, err := mm.CurrentExecApp()
			if err != nil {
				return mm, nil
			}
			return mm, mm.RunExecApp(currentExecApp)
		case key.Matches(msg, mm.keys.Stop):
			currentExecApp, err := mm.CurrentExecApp()
			if err != nil {
				return mm, nil
			}
			return mm, mm.StopExecApp(currentExecApp)
		case key.Matches(msg, mm.keys.RunAll):
			m, cmd := mm.checkRequiredVar(func() tea.Msg { return msg })
			if m != nil {
				return m, cmd
			}
			cmds := make([]tea.Cmd, 0)
			for _, execApp := range mm.ExecutableApps {
				if cmd := mm.RunExecApp(execApp); cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
			return mm, tea.Batch(cmds...)
		case key.Matches(msg, mm.keys.StopAll):
			cmds := make([]tea.Cmd, 0)
			for _, execApp := range mm.ExecutableApps {
				if cmd := mm.StopExecApp(execApp); cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
			return mm, tea.Batch(cmds...)
		case key.Matches(msg, mm.keys.Vars):
			return mm.varsModel, nil
		case key.Matches(msg, mm.keys.Check):
			cmds := make([]tea.Cmd, 0)
			for _, execApp := range mm.ExecutableApps {
				if cmd := mm.CheckExecApp(execApp, 0); cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
			return mm, tea.Batch(cmds...)

		}
	case UpdateExecAppMsg:
		return mm, checkExecutableApp(msg)
	case tea.WindowSizeMsg:
		mm.height = msg.Height
		mm.width = msg.Width
		mm.help.Width = msg.Width
		mm.varsModel.Height = msg.Height
	case TickMsg:
		return mm, mm.doTick()
	}

	return mm, nil
}

type UpdateExecAppMsg *executable.ExecutableApp

func checkExecutableApp(ea *executable.ExecutableApp) tea.Cmd {
	if ((ea.Status() == executable.Stopping || ea.Status() == executable.WithError) && ea.PreviousStatus() == executable.Running) || (ea.Status() == executable.IsNotRunning && ea.PreviousStatus() == executable.Stopping) {
		return nil
	} else {
		return func() tea.Msg {
			return UpdateExecAppMsg(ea)
		}
	}
}
