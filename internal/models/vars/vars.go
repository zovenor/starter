package vars

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/zovenor/starter/internal/config"
	"github.com/zovenor/starter/internal/keymap"
)

type VarWithValue struct {
	config.Var
	Value string
}

func GetValueByName(name string, vrs []*VarWithValue) (string, error) {
	for _, vwv := range vrs {
		if vwv.Name == name {
			return vwv.Value, nil
		}
	}
	return "", fmt.Errorf("can not find var with name %v", name)
}

type VarsModel struct {
	Vars           []*VarWithValue
	parent         tea.Model
	cursor         int
	keys           keymap.KeyMap
	help           help.Model
	Height         int
	parentName     string
	updatedValue   int
	textInputModel *textinput.Model
}

func New(varsConfig []config.Var, parent tea.Model, keys keymap.KeyMap, helpModel help.Model, parentName string) *VarsModel {
	vars := make([]*VarWithValue, 0)
	for _, varCfg := range varsConfig {
		vars = append(vars, &VarWithValue{
			Var: varCfg,
		})
	}
	vm := new(VarsModel)
	vm.Vars = vars
	vm.parent = parent
	vm.keys = keys
	vm.help = helpModel
	vm.parentName = parentName
	return vm
}

func (vm *VarsModel) Init() tea.Cmd {
	return nil
}

func (vm *VarsModel) View() string {
	s := fmt.Sprintf("%v • %v\n\n", vm.parentName, "Vars")

	for i, v := range vm.Vars {
		if i == vm.cursor {
			s += "‣ "
		} else {
			s += "  "
		}
		s += fmt.Sprintf("%v", v.Name)
		if v.Required {
			s += color.New(color.FgYellow).Sprint(" (required)")
		}
		s += ": "
		if vm.updatedValue == i && vm.textInputModel != nil {
			s += color.New(color.Bold).Sprintf("%v", vm.textInputModel.View())
		} else {
			if v.Hiden {
				s += strings.Repeat("*", len(v.Value))
			} else {
				s += v.Value
			}
		}
		s += "\n"
	}

	keysView := vm.help.View(vm.keys)
	keysView += "\n"
	rpt := vm.Height - strings.Count(s, "\n") - strings.Count(keysView, "\n") - 1
	if rpt < 0 {
		rpt = 0
	}
	s += strings.Repeat("\n", rpt)
	s += keysView

	return s
}

func (vm *VarsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if vm.textInputModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, vm.keys.Cancel):
				vm.textInputModel = nil
				vm.updatedValue = -1
				return vm, nil
			case key.Matches(msg, vm.keys.Enter):
				value := vm.textInputModel.Value()
				vm.Vars[vm.updatedValue].Value = value
				vm.textInputModel = nil
				vm.updatedValue = -1
				return vm, nil
			}
		}
		newTI, cmd := vm.textInputModel.Update(msg)
		vm.textInputModel = &newTI
		return vm, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, vm.keys.Quit):
			return vm, tea.Quit
		case key.Matches(msg, vm.keys.Cancel):
			return vm.parent, nil
		case key.Matches(msg, vm.keys.Up):
			if vm.cursor > 0 {
				vm.cursor--
			}
		case key.Matches(msg, vm.keys.Down):
			if vm.cursor < len(vm.Vars)-1 {
				vm.cursor++
			}
		case key.Matches(msg, vm.keys.Edit):
			ti := textinput.New()
			ti.SetValue(vm.Vars[vm.cursor].Value)
			ti.Placeholder = vm.Vars[vm.cursor].Value
			ti.Cursor.SetMode(cursor.CursorHide)
			ti.Focus()
			vm.textInputModel = &ti
			vm.updatedValue = vm.cursor
			return vm, nil
		}
	}
	return vm, nil
}
