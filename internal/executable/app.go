package executable

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/zovenor/starter/internal/config"
)

type Status uint8

const (
	IsNotRunning Status = iota
	Running
	Executed
	WithError
	Stopping
)

type ExecutableApp struct {
	config.ExecAppConfig

	Index int
	Log   string

	status         Status
	previousStatus Status
}

func (execApp *ExecutableApp) SetStatus(status Status) {
	execApp.previousStatus = execApp.status
	execApp.status = status
}

func (execApp *ExecutableApp) Status() Status {
	return execApp.status
}

func (execApp *ExecutableApp) PreviousStatus() Status {
	return execApp.previousStatus
}

func (execApp *ExecutableApp) Format(selected bool) string {
	itemLine := fmt.Sprintf("○ %v", execApp.Name)
	if execApp.Log != "" {
		itemLine += fmt.Sprintf(": %v", execApp.Log)
	}

	if !execApp.Disabled {
		switch execApp.Status() {
		case Running:
			itemLine = color.New(color.FgBlue).Sprint(itemLine)
		case Executed:
			itemLine = color.New(color.FgGreen).Sprint(itemLine)
		case WithError:
			itemLine = color.New(color.FgRed).Sprint(itemLine)
		}
	} else {
		itemLine = color.New(color.FgHiBlack).Sprint(itemLine)
	}
	if selected {
		itemLine = fmt.Sprintf(" > %v", itemLine)
	} else {
		itemLine = fmt.Sprintf("   %v", itemLine)
	}
	itemLine += "\n"
	return itemLine
}

func (execApp *ExecutableApp) Run() {
	execApp.Log = "Running..."
	execApp.SetStatus(Running)
	err := execApp.runCmds()
	if err != nil {
		execApp.Log = err.Error()
		execApp.SetStatus(WithError)
	} else {
		execApp.Log = "Executed"
		execApp.SetStatus(Executed)
	}
}

func (execApp *ExecutableApp) Stop() {
	execApp.Log = "Stopping..."
	execApp.SetStatus(Stopping)
	err := execApp.stopCmds()
	if err != nil {
		execApp.Log = err.Error()
	}
	execApp.Log = ""
	execApp.SetStatus(IsNotRunning)
}

func (execApp *ExecutableApp) runCmds() error {
	for _, cmdString := range execApp.Cmds {
		execApp.Log = cmdString
		time.Sleep(time.Second)
		cList := strings.Split(cmdString, " ")
		if len(cList) == 0 {
			return fmt.Errorf("len of command is zero")
		}
		cAppString := cList[0]
		args := make([]string, 0)
		if len(cList) > 1 {
			args = cList[1:]
		}
		cmd := exec.Command(cAppString, args...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (execApp *ExecutableApp) stopCmds() error {
	for _, cmdString := range execApp.StopCmds {
		execApp.Log = cmdString
		time.Sleep(time.Second)
		cList := strings.Split(cmdString, " ")
		if len(cList) == 0 {
			return fmt.Errorf("len of command is zero")
		}
		cAppString := cList[0]
		args := make([]string, 0)
		if len(cList) > 1 {
			args = cList[1:]
		}
		cmd := exec.Command(cAppString, args...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
