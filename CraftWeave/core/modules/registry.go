package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

type TaskHandler interface {
	Run(h inventory.Host, task parser.Task, vars map[string]string) ssh.CommandResult
}

var registry = make(map[string]TaskHandler)

func Register(name string, h TaskHandler) {
	registry[name] = h
}

func GetHandler(name string) TaskHandler {
	return registry[name]
}
