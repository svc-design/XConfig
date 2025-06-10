// core/parser/parser.go
package parser

type Template struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type Task struct {
	Name     string    `yaml:"name"`
	When     string    `yaml:"when,omitempty"`
	Shell    string    `yaml:"shell,omitempty"`
	Script   string    `yaml:"script,omitempty"`
	Template *Template `yaml:"template,omitempty"`
}

// Type returns the module name associated with this task.
func (t Task) Type() string {
	switch {
	case t.Shell != "":
		return "shell"
	case t.Script != "":
		return "script"
	case t.Template != nil:
		return "template"
	default:
		return ""
	}
}

type Play struct {
	Name  string            `yaml:"name"`
	Hosts string            `yaml:"hosts"`
	Vars  map[string]string `yaml:"vars,omitempty"`
	Tasks []Task            `yaml:"tasks"`
}
