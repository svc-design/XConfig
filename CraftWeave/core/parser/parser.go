// core/parser/parser.go
package parser

type Template struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type Task struct {
	Name     string    `yaml:"name"`
	Shell    string    `yaml:"shell,omitempty"`
	Script   string    `yaml:"script,omitempty"`
	Template *Template `yaml:"template,omitempty"`
}

type Play struct {
	Name  string            `yaml:"name"`
	Hosts string            `yaml:"hosts"`
	Vars  map[string]string `yaml:"vars,omitempty"`
	Tasks []Task            `yaml:"tasks"`
}
