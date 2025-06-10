package modules

import (
	"bytes"
	"text/template"

	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

// ShellHandler executes a shell command on remote host
type ShellHandler struct{}

func (ShellHandler) Run(h inventory.Host, task parser.Task, vars map[string]string) ssh.CommandResult {
	cmd := task.Shell
	if len(vars) > 0 {
		t, err := template.New("shell").Parse(cmd)
		if err == nil {
			var buf bytes.Buffer
			if err := t.Execute(&buf, vars); err == nil {
				cmd = buf.String()
			}
		}
	}
	return ssh.RunShellCommand(h, cmd)
}

func init() {
	Register("shell", ShellHandler{})
}
