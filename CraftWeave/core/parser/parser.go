// core/parser/parser.go
package parser

type Template struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

// Stat describes the stat module options
type Stat struct {
	Path string `yaml:"path"`
}

// Fail describes the fail module options
type Fail struct {
	Msg string `yaml:"msg"`
}

// Debug describes the debug module options
type Debug struct {
	Msg string `yaml:"msg"`
}

type Task struct {
	Name     string    `yaml:"name"`
	Shell    string    `yaml:"shell,omitempty"`
	Script   string    `yaml:"script,omitempty"`
	Template *Template `yaml:"template,omitempty"`
	Stat     *Stat     `yaml:"stat,omitempty"`
	Fail     *Fail     `yaml:"fail,omitempty"`
	Debug    *Debug    `yaml:"debug,omitempty"`
	Register string    `yaml:"register,omitempty"`
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
