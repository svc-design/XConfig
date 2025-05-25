// core/parser/parser.go
package parser

type Task struct {
	Name   string `yaml:"name"`
	Shell  string `yaml:"shell,omitempty"`
	Script string `yaml:"script,omitempty"`
}

type Play struct {
	Name   string `yaml:"name"`
	Hosts  string `yaml:"hosts"`
	Tasks  []Task `yaml:"tasks"`
}
