package executable

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/zovenor/starter/internal/config"
	"github.com/zovenor/starter/internal/models/vars"
)

type Status uint8

const (
	IsNotRunning Status = iota
	Running
	Executed
	WithError
	Stopping
)

const CheckingLog = "Checking..."

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
		case Stopping:
			itemLine = color.New(color.FgYellow).Sprint(itemLine)
		}
	} else {
		itemLine = color.New(color.FgHiBlack).Sprint(itemLine)
	}
	if selected {
		itemLine = fmt.Sprintf(" ‣ %v", itemLine)
	} else {
		itemLine = fmt.Sprintf("   %v", itemLine)
	}
	itemLine += "\n"
	return itemLine
}

func (execApp *ExecutableApp) Run(vrs []*vars.VarWithValue) {
	if execApp.Disabled {
		return
	}
	execApp.Stop()
	execApp.Log = "Running..."
	execApp.SetStatus(Running)
	err := execApp.runCmds(vrs)
	if err != nil {
		execApp.Log = err.Error()
		execApp.SetStatus(WithError)
	} else {
		execApp.Log = "Executed"
		execApp.SetStatus(Executed)
	}
}

func (execApp *ExecutableApp) Stop() {
	if execApp.Disabled {
		return
	}
	execApp.Log = "Stopping..."
	execApp.SetStatus(Stopping)
	err := execApp.stopCmds()
	if err != nil {
		execApp.Log = err.Error()
	}
	execApp.Log = ""
	execApp.SetStatus(IsNotRunning)
}

func (execApp *ExecutableApp) runCmds(vrs []*vars.VarWithValue) error {
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("error related to get currect user: %v", err.Error())
	}
	for _, cmdString := range execApp.Cmds {
		execApp.Log = cmdString
		time.Sleep(time.Second)
		newCmdString := cmdString
		newCmdString = strings.ReplaceAll(newCmdString, "~", user.HomeDir)
		newCmdString = strings.ReplaceAll(newCmdString, "$HOME", user.HomeDir)
		cList := strings.Fields(newCmdString)
		if len(cList) == 0 {
			return fmt.Errorf("len of command is zero")
		}
		args := make([]string, 0)
		if len(cList) > 1 {
			args = cList[1:]
		}
		cmd := exec.Command(cList[0], args...)

		lastEnv := cmd.Env
		cmd.Env = append(cmd.Env, os.Environ()...)
		for _, v := range vrs {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", v.Name, v.Value))
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			cmd.Env = lastEnv
			return fmt.Errorf("%v (cmd: %v, output: %v)", err, newCmdString, string(output))
		}
		cmd.Env = lastEnv
		execApp.Log = string(output)
	}
	return nil
}

func (execApp *ExecutableApp) stopCmds() error {
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("error related to get currect user: %v", err.Error())
	}
	for _, cmdString := range execApp.StopCmds {
		execApp.Log = cmdString
		time.Sleep(time.Second)
		newCmdString := cmdString
		newCmdString = strings.ReplaceAll(newCmdString, "~", user.HomeDir)
		newCmdString = strings.ReplaceAll(newCmdString, "$HOME", user.HomeDir)
		cList := strings.Fields(newCmdString)
		if len(cList) == 0 {
			return fmt.Errorf("len of command is zero")
		}
		args := make([]string, 0)
		if len(cList) > 1 {
			args = cList[1:]
		}
		cmd := exec.Command(cList[0], args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v (cmd: %v, output: %v)", err, newCmdString, string(output))
		}
		execApp.Log = string(output)
	}
	return nil
}

func (execApp *ExecutableApp) Check() error {
	if execApp.Disabled {
		return nil
	}
	if execApp.Log == CheckingLog {
		return nil
	}
	if execApp.CheckCmd == "" {
		return nil
	}
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("error related to get currect user: %v", err.Error())
	}
	execApp.SetStatus(Running)
	execApp.Log = CheckingLog
	time.Sleep(time.Second)
	newCmdString := execApp.CheckCmd
	newCmdString = strings.ReplaceAll(newCmdString, "~", user.HomeDir)
	newCmdString = strings.ReplaceAll(newCmdString, "$HOME", user.HomeDir)
	cList := strings.Fields(newCmdString)
	if len(cList) == 0 {
		return fmt.Errorf("len of command is zero")
	}
	args := make([]string, 0)
	if len(cList) > 1 {
		args = cList[1:]
	}
	cmd := exec.Command(cList[0], args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		execApp.Log = fmt.Sprintf("Output: %s, Error: %s", output, err.Error())
		execApp.SetStatus(WithError)
	} else {
		execApp.Log = "Executed"
		execApp.SetStatus(Executed)
	}
	return nil
}
