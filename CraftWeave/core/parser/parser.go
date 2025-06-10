// core/parser/parser.go
package parser

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

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

type CommandTask struct {
	Cmd string `yaml:"cmd,omitempty"`
}

type AptTask struct {
	Name  string `yaml:"name,omitempty"`
	Deb   string `yaml:"deb,omitempty"`
	State string `yaml:"state,omitempty"`
}

type YumTask struct {
	Name  string `yaml:"name,omitempty"`
	State string `yaml:"state,omitempty"`
}

type SystemdTask struct {
	Name    string `yaml:"name"`
	State   string `yaml:"state,omitempty"`
	Enabled bool   `yaml:"enabled,omitempty"`
}

type ServiceTask struct {
	Name    string `yaml:"name"`
	State   string `yaml:"state,omitempty"`
	Enabled bool   `yaml:"enabled,omitempty"`
}

type FailTask struct {
	Msg string `yaml:"msg"`
}

type DebugTask struct {
	Msg string `yaml:"msg"`
}

type StatTask struct {
	Path string `yaml:"path"`
}

type SetFactTask map[string]string

type SetupTask struct{}

type Task struct {
	Name     string       `yaml:"name"`
	Shell    string       `yaml:"shell,omitempty"`
	Script   string       `yaml:"script,omitempty"`
	Template *Template    `yaml:"template,omitempty"`
	Copy     *CopyTask    `yaml:"copy,omitempty"`
	Command  string       `yaml:"command,omitempty"`
	Apt      *AptTask     `yaml:"apt,omitempty"`
	Yum      *YumTask     `yaml:"yum,omitempty"`
	Systemd  *SystemdTask `yaml:"systemd,omitempty"`
	Service  *ServiceTask `yaml:"service,omitempty"`
	Fail     *FailTask    `yaml:"fail,omitempty"`
	Debug    *DebugTask   `yaml:"debug,omitempty"`
	Stat     *StatTask    `yaml:"stat,omitempty"`
	SetFact  SetFactTask  `yaml:"set_fact,omitempty"`
	Setup    *SetupTask   `yaml:"setup,omitempty"`
}

// Package defines a simple package installation task
type Package struct {
	Name  string `yaml:"name"`
	State string `yaml:"state,omitempty"`
}

// Role reference used in Play definition
type Role struct {
	Role string `yaml:"role"`
}

type Play struct {
	Name        string            `yaml:"name"`
	Hosts       string            `yaml:"hosts"`
	GatherFacts *bool             `yaml:"gather_facts,omitempty"`
	Vars        map[string]string `yaml:"vars,omitempty"`
	Tasks       []Task            `yaml:"tasks"`
	Roles       []Role            `yaml:"roles,omitempty"`
}

// UnmarshalYAML ensures only known task fields are accepted
func (t *Task) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]*yaml.Node
	if err := value.Decode(&raw); err != nil {
		return err
	}

	allowed := map[string]bool{
		"name":     true,
		"shell":    true,
		"script":   true,
		"template": true,
		"copy":     true,
		"command":  true,
		"apt":      true,
		"yum":      true,
		"systemd":  true,
		"service":  true,
		"fail":     true,
		"debug":    true,
		"stat":     true,
		"set_fact": true,
		"setup":    true,
	}

	for k := range raw {
		if !allowed[k] {
			return fmt.Errorf("unknown task field '%s'", k)
		}
	}

	type alias Task
	var a alias
	if err := value.Decode(&a); err != nil {
		return err
	}
	*t = Task(a)
	return nil
}
