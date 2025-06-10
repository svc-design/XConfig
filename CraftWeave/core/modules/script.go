package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

type ScriptHandler struct{}

func (ScriptHandler) Run(h inventory.Host, task parser.Task, vars map[string]string) ssh.CommandResult {
	return ssh.RunRemoteScript(h, task.Script)
}

func init() {
	Register("script", ScriptHandler{})
}
