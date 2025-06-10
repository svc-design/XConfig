// core/parser/parser.go
package parser

type Template struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

// Copy defines parameters for the copy/rsync module
type Copy struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
	Mode string `yaml:"mode,omitempty"`
}

type Task struct {
	Name     string    `yaml:"name"`
	Shell    string    `yaml:"shell,omitempty"`
	Script   string    `yaml:"script,omitempty"`
	Template *Template `yaml:"template,omitempty"`
	Copy     *Copy     `yaml:"copy,omitempty"`
}

// Role reference used in Play definition
type Role struct {
	Role string `yaml:"role"`
}

type Play struct {
	Name  string            `yaml:"name"`
	Hosts string            `yaml:"hosts"`
	Vars  map[string]string `yaml:"vars,omitempty"`
	Tasks []Task            `yaml:"tasks"`
	Roles []Role            `yaml:"roles,omitempty"`
}
