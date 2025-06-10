package inventory

import (
	"bufio"
	"os"
	"strings"
)

type Host struct {
	Name     string
	Address  string
	User     string
	KeyFile  string
	Port     string
	Password string            // ✅ 新增：支持密码登录
	Vars     map[string]string // ✅ 新增：主机变量
}

func Parse(path, group string) ([]Host, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// group -> []*Host for later variable merge
	groupHosts := make(map[string][]*Host)
	groupVars := make(map[string]map[string]string)

	scanner := bufio.NewScanner(file)
	currentGroup := ""
	inVars := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.Trim(line, "[]")
			if strings.HasSuffix(section, ":vars") {
				currentGroup = strings.TrimSuffix(section, ":vars")
				inVars = true
			} else {
				currentGroup = section
				inVars = false
			}
			continue
		}

		if inVars {
			kv := strings.SplitN(line, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				val := strings.Trim(strings.TrimSpace(kv[1]), "\"'")
				if _, ok := groupVars[currentGroup]; !ok {
					groupVars[currentGroup] = make(map[string]string)
				}
				groupVars[currentGroup][key] = val
			}
			continue
		}

		// host definitions
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		h := Host{Name: parts[0], Address: parts[0], Vars: make(map[string]string)}
		for _, p := range parts[1:] {
			if !strings.Contains(p, "=") {
				continue
			}
			kv := strings.SplitN(p, "=", 2)
			key := kv[0]
			val := strings.Trim(kv[1], "\"'")

			switch key {
			case "ansible_host":
				h.Address = val
			case "ansible_ssh_user":
				h.User = val
			case "ansible_port":
				h.Port = val
			case "ansible_ssh_private_key_file":
				h.KeyFile = val
			case "ansible_ssh_pass":
				h.Password = val
			default:
				h.Vars[key] = val
			}
		}
		groupHosts[currentGroup] = append(groupHosts[currentGroup], &h)
	}

	// defaults from [all:vars]
	defaults := map[string]string{
		"ansible_ssh_user":             "ubuntu",
		"ansible_port":                 "22",
		"ansible_ssh_private_key_file": os.Getenv("HOME") + "/.ssh/id_rsa",
	}
	if gv, ok := groupVars["all"]; ok {
		for k, v := range gv {
			defaults[k] = v
		}
	}

	var hosts []Host
	for _, hptr := range groupHosts[group] {
		h := hptr
		// merge connection info from group vars and defaults
		if h.User == "" {
			if v, ok := groupVars[group]["ansible_ssh_user"]; ok {
				h.User = v
			} else {
				h.User = defaults["ansible_ssh_user"]
			}
		}
		if h.Port == "" {
			if v, ok := groupVars[group]["ansible_port"]; ok {
				h.Port = v
			} else {
				h.Port = defaults["ansible_port"]
			}
		}
		if h.KeyFile == "" {
			if v, ok := groupVars[group]["ansible_ssh_private_key_file"]; ok {
				h.KeyFile = v
			} else {
				h.KeyFile = defaults["ansible_ssh_private_key_file"]
			}
		}
		if h.Password == "" {
			if v, ok := groupVars[group]["ansible_ssh_pass"]; ok {
				h.Password = v
			} else if v, ok := defaults["ansible_ssh_pass"]; ok {
				h.Password = v
			}
		}

		// merge group vars and defaults into host vars (special keys excluded)
		mergeVars := func(vars map[string]string) {
			for k, v := range vars {
				switch k {
				case "ansible_host", "ansible_ssh_user", "ansible_port", "ansible_ssh_private_key_file", "ansible_ssh_pass":
					continue
				}
				if _, ok := h.Vars[k]; !ok {
					h.Vars[k] = v
				}
			}
		}
		if gv, ok := groupVars["all"]; ok {
			mergeVars(gv)
		}
		if gv, ok := groupVars[group]; ok {
			mergeVars(gv)
		}

		hosts = append(hosts, *h)
	}

	return hosts, nil
}
