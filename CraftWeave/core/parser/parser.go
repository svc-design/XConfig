// core/parser/parser.go
package parser

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
	Roles []struct {
		Role string `yaml:"role"`
	} `yaml:"roles,omitempty"`
	Tasks []Task `yaml:"tasks,omitempty"`
}

// LoadPlaybook parses the given playbook YAML and expands any referenced roles.
func LoadPlaybook(path string) ([]Play, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var plays []Play
	if err := yaml.Unmarshal(data, &plays); err != nil {
		return nil, err
	}

	base := filepath.Dir(path)
	for i := range plays {
		var allTasks []Task
		for _, r := range plays[i].Roles {
			ts, err := loadRoleTasks(base, r.Role)
			if err != nil {
				return nil, err
			}
			allTasks = append(allTasks, ts...)
		}
		allTasks = append(allTasks, plays[i].Tasks...)
		plays[i].Tasks = allTasks
	}

	return plays, nil
}

func loadRoleTasks(base, name string) ([]Task, error) {
	dir := filepath.Join(base, "roles", name, "tasks")
	path := filepath.Join(dir, "main.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(dir, "main.yml")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := yaml.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
