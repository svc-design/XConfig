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
	Systemd  *Service  `yaml:"systemd,omitempty"`
	Service  *Service  `yaml:"service,omitempty"`
}

type Service struct {
	Name    string `yaml:"name"`
	State   string `yaml:"state,omitempty"`
	Enabled *bool  `yaml:"enabled,omitempty"`
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
