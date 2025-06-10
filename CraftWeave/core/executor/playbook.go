// core/executor/playbook.go
package executor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"gopkg.in/yaml.v3"

	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

// 全局控制输出样式和 dry-run 模式
var AggregateOutput bool
var CheckMode bool

// ExecutePlaybook 解析并执行整个 playbook
func ExecutePlaybook(playbook []parser.Play, inventoryPath string, baseDir string, extraVars map[string]string) {
	for _, play := range playbook {
		fmt.Printf("\n🎯 Play: %s (hosts: %s)\n", play.Name, play.Hosts)

		hosts, err := inventory.Parse(inventoryPath, play.Hosts)
		if err != nil {
			fmt.Printf("❌ Failed to resolve hosts: %v\n", err)
			continue
		}

		allTasks := append([]parser.Task{}, play.Tasks...)
		for _, role := range play.Roles {
			roleDir := filepath.Join(baseDir, "roles", role.Role)
			rolePath := filepath.Join(roleDir, "tasks", "main.yaml")
			data, err := os.ReadFile(rolePath)
			if err != nil {
				fmt.Printf("❌ Failed to load role %s: %v\n", role.Role, err)
				continue
			}
			var roleTasks []parser.Task
			if err := yaml.Unmarshal(data, &roleTasks); err != nil {
				fmt.Printf("❌ Failed to parse role %s: %v\n", role.Role, err)
				continue
			}

			for i := range roleTasks {
				if roleTasks[i].Script != "" && !filepath.IsAbs(roleTasks[i].Script) {
					scriptPath := filepath.Join(roleDir, roleTasks[i].Script)
					if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
						alt := filepath.Join(roleDir, "scripts", roleTasks[i].Script)
						if _, err := os.Stat(alt); err == nil {
							scriptPath = alt
						}
					}
					roleTasks[i].Script = scriptPath
				}
				if roleTasks[i].Template != nil && roleTasks[i].Template.Src != "" && !filepath.IsAbs(roleTasks[i].Template.Src) {
					tplPath := filepath.Join(roleDir, roleTasks[i].Template.Src)
					if _, err := os.Stat(tplPath); os.IsNotExist(err) {
						alt := filepath.Join(roleDir, "templates", roleTasks[i].Template.Src)
						if _, err := os.Stat(alt); err == nil {
							tplPath = alt
						}
					}
					roleTasks[i].Template.Src = tplPath
				}
			}
			allTasks = append(allTasks, roleTasks...)
		}

		var results []ssh.CommandResult
		var mu sync.Mutex
		var wg sync.WaitGroup

		// merge play vars with extra vars (extra vars override)
		playVars := make(map[string]string)
		for k, v := range play.Vars {
			playVars[k] = v
		}
		for k, v := range extraVars {
			playVars[k] = v
		}

		hostFacts := make(map[string]map[string]string)
		gather := true
		if play.GatherFacts != nil {
			gather = *play.GatherFacts
		}
		if gather {
			for _, h := range hosts {
				hostFacts[h.Name] = ssh.GatherFacts(h)
			}
		}

		varsPerHost := make(map[string]map[string]string)
		for _, h := range hosts {
			hv := make(map[string]string)
			// seed with gathered facts for the host
			if f, ok := hostFacts[h.Name]; ok {
				for k, v := range f {
					hv[k] = v
				}
			}
			// play level variables (merged with extra vars above)
			for k, v := range playVars {
				hv[k] = v
			}
			// host specific variables from inventory
			for k, v := range h.Vars {
				hv[k] = v
			}
			varsPerHost[h.Name] = hv
		}

		for _, host := range hosts {
			wg.Add(1)
			go func(h inventory.Host) {
				defer wg.Done()

				hv := varsPerHost[h.Name]

				for _, task := range allTasks {

					if CheckMode {
						fmt.Printf("%s | SKIPPED | dry-run: %s\n", h.Name, task.Name)
						continue
					}

					if task.When != "" && !evaluateCondition(task.When, hv) {
						mu.Lock()
						results = append(results, ssh.CommandResult{
							Host:       h.Name,
							ReturnMsg:  "SKIPPED",
							ReturnCode: 0,
							Output:     fmt.Sprintf("when condition '%s' not met", task.When),
						})
						mu.Unlock()
						continue
					}

					var res ssh.CommandResult

					if task.Command != "" {
						rendered := task.Command
						if len(hv) > 0 {
							renderedTmpl, err := template.New("command").Parse(task.Command)
							if err == nil {
								var buf bytes.Buffer
								if err := renderedTmpl.Execute(&buf, hv); err == nil {
									rendered = buf.String()
								}
							}
						}
						res = ssh.RunCommand(h, rendered)
					} else if task.Shell != "" {
						rendered := task.Shell
						if len(hv) > 0 {
							renderedTmpl, err := template.New("shell").Parse(task.Shell)
							if err == nil {
								var buf bytes.Buffer
								if err := renderedTmpl.Execute(&buf, hv); err == nil {
									rendered = buf.String()
								}
							}
						}
						res = ssh.RunShellCommand(h, rendered)
					} else if task.Script != "" {
						res = ssh.RunRemoteScript(h, task.Script)
					} else if task.Template != nil {
						res = ssh.RenderTemplate(h, task.Template.Src, task.Template.Dest, hv)
					} else if len(task.SetFact) > 0 {
						for k, v := range task.SetFact {
							val := v
							if tmpl, err := template.New("sf").Parse(v); err == nil {
								var buf bytes.Buffer
								if err := tmpl.Execute(&buf, hv); err == nil {
									val = buf.String()
								}
							}
							hv[k] = val
						}
						res = ssh.RenderTemplate(h, task.Template.Src, task.Template.Dest, hv)
					} else if task.Copy != nil {
						src := task.Copy.Src
						dest := task.Copy.Dest
						if len(hv) > 0 {
							// render src and dest with variables if needed
							for _, field := range []struct {
								val *string
							}{
								{&src}, {&dest},
							} {
								tmpl, err := template.New("copy").Parse(*field.val)
								if err == nil {
									var buf bytes.Buffer
									if err := tmpl.Execute(&buf, hv); err == nil {
										*field.val = buf.String()
									}
								}
							}
						}
						res = ssh.CopyFile(h, src, dest, task.Copy.Mode)
					} else if task.Apt != nil {
						if task.Apt.State == "" || task.Apt.State == "present" {
							res = ssh.InstallAptPackage(h, task.Apt.Name)
						} else {
							res = ssh.CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("unsupported state '%s'", task.Apt.State)}
						}
					} else if task.Yum != nil {
						if task.Yum.State == "" || task.Yum.State == "present" {
							res = ssh.InstallYumPackage(h, task.Yum.Name)
						} else {
							res = ssh.CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("unsupported state '%s'", task.Yum.State)}
						}
					} else {
						res = ssh.CommandResult{
							Host:       h.Name,
							ReturnMsg:  "FAILED",
							ReturnCode: 1,
							Output:     fmt.Sprintf("Unsupported task type in '%s'", task.Name),
						}
					}

					if task.Register != "" {
						hv[task.Register] = strings.TrimSpace(res.Output)
					}

					if len(task.SetFact) > 0 && (task.Shell != "" || task.Script != "" || task.Template != nil) {
						for k, v := range task.SetFact {
							val := v
							if tmpl, err := template.New("sf").Parse(v); err == nil {
								var buf bytes.Buffer
								if err := tmpl.Execute(&buf, hv); err == nil {
									val = buf.String()
								}
							}
							hv[k] = val
						}
					}

					mu.Lock()
					results = append(results, res)
					mu.Unlock()
				}
			}(host)
		}
		wg.Wait()

		if AggregateOutput {
			ssh.AggregatedPrint(results)
		} else {
			for _, r := range results {
				fmt.Printf("%s | %s | rc=%d >>\n%s\n", r.Host, r.ReturnMsg, r.ReturnCode, r.Output)
			}
		}
	}
}

func evaluateCondition(cond string, vars map[string]string) bool {
	cond = strings.TrimSpace(cond)
	if cond == "" {
		return true
	}
	if strings.Contains(cond, "==") {
		parts := strings.SplitN(cond, "==", 2)
		left := strings.TrimSpace(parts[0])
		right := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		return vars[left] == right
	}
	if strings.Contains(cond, "!=") {
		parts := strings.SplitN(cond, "!=", 2)
		left := strings.TrimSpace(parts[0])
		right := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		return vars[left] != right
	}
	val, ok := vars[cond]
	return ok && val != "" && val != "false"
}
