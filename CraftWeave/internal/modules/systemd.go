package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"fmt"
)

func systemdHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Systemd == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing systemd parameters"}
	}
	cmd := fmt.Sprintf("sudo systemctl %s %s", task.Systemd.State, task.Systemd.Name)
	if task.Systemd.Enabled {
		cmd = fmt.Sprintf("%s && sudo systemctl enable %s", cmd, task.Systemd.Name)
	}
	return ssh.RunShellCommand(ctx.Host, cmd)
}

func init() { Register("systemd", systemdHandler) }
