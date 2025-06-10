package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

type TemplateHandler struct{}

func (TemplateHandler) Run(h inventory.Host, task parser.Task, vars map[string]string) ssh.CommandResult {
	if task.Template == nil {
		return ssh.CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "template data missing"}
	}
	return ssh.RenderTemplate(h, task.Template.Src, task.Template.Dest, vars)
}

func init() {
	Register("template", TemplateHandler{})
}
