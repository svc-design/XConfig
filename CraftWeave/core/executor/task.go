package executor

import (
	"fmt"

	"craftweave/core/modules"
	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

// ExecuteTask executes a task using registered modules
func ExecuteTask(h inventory.Host, task parser.Task, vars map[string]string) ssh.CommandResult {
	handler := modules.GetHandler(task.Type())
	if handler == nil {
		return ssh.CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("unsupported module '%s'", task.Type()),
		}
	}
	return handler.Run(h, task, vars)
}
