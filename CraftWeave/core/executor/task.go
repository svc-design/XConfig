package executor

import (
	"fmt"

	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/modules"
	"craftweave/internal/ssh"
)

// ExecuteTask dispatches the task to the appropriate module handler.
func ExecuteTask(task parser.Task, host inventory.Host, vars map[string]string) ssh.CommandResult {
	ctx := modules.Context{Host: host, Vars: vars}

	if h, ok := modules.GetHandler(task.Type()); ok {
		return h(ctx, task)
	}

	switch {
	case task.Shell != "":
		return ssh.RunShellCommand(host, task.Shell)
	case task.Script != "":
		return ssh.RunRemoteScript(host, task.Script)
	case task.Template != nil:
		return ssh.RenderTemplate(host, task.Template.Src, task.Template.Dest, vars)
	default:
		return ssh.CommandResult{Host: host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("Unsupported task type in '%s'", task.Name)}
	}
}
