package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"fmt"
)

func serviceHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Service == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing service parameters"}
	}
	cmd := fmt.Sprintf("sudo service %s %s", task.Service.Name, task.Service.State)
	if task.Service.Enabled {
		cmd = fmt.Sprintf("%s && sudo systemctl enable %s", cmd, task.Service.Name)
	}
	return ssh.RunShellCommand(ctx.Host, cmd)
}

func init() { Register("service", serviceHandler) }
