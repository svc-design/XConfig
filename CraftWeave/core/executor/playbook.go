// core/executor/playbook.go
package executor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
			rolePath := filepath.Join(baseDir, "roles", role.Role, "tasks", "main.yaml")
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
					roleTasks[i].Script = filepath.Join(baseDir, "roles", role.Role, roleTasks[i].Script)
				}
				if roleTasks[i].Template != nil && roleTasks[i].Template.Src != "" && !filepath.IsAbs(roleTasks[i].Template.Src) {
					roleTasks[i].Template.Src = filepath.Join(baseDir, "roles", role.Role, roleTasks[i].Template.Src)
				}
			}
			allTasks = append(allTasks, roleTasks...)
		}

		var results []ssh.CommandResult
		var mu sync.Mutex
		var wg sync.WaitGroup

		// merge play vars with extra vars (extra vars override)
		playVars := make(map[string]interface{})
		for k, v := range play.Vars {
			playVars[k] = v
		}
		for k, v := range extraVars {
			playVars[k] = v
		}

		for _, host := range hosts {
			host := host
			wg.Add(1)
			go func(h inventory.Host) {
				defer wg.Done()

				hostVars := make(map[string]interface{})
				for k, v := range playVars {
					hostVars[k] = v
				}

				for _, task := range allTasks {
					if CheckMode {
						fmt.Printf("%s | SKIPPED | dry-run: %s\n", h.Name, task.Name)
						continue
					}

					var res ssh.CommandResult

					// Handle shell module
					if task.Shell != "" {
						rendered := task.Shell
						if len(hostVars) > 0 {
							renderedTmpl, err := template.New("shell").Parse(task.Shell)
							if err == nil {
								var buf bytes.Buffer
								if err := renderedTmpl.Execute(&buf, hostVars); err == nil {
									rendered = buf.String()
								}
							}
						}
						res = ssh.RunShellCommand(h, rendered)
					} else if task.Script != "" { // script module
						res = ssh.RunRemoteScript(h, task.Script)
					} else if task.Template != nil {
						res = ssh.RenderTemplate(h, task.Template.Src, task.Template.Dest, hostVars)
					} else if task.Stat != nil { // stat module
						path := task.Stat.Path
						if len(hostVars) > 0 {
							tmpl, err := template.New("stat").Parse(task.Stat.Path)
							if err == nil {
								var buf bytes.Buffer
								if err := tmpl.Execute(&buf, hostVars); err == nil {
									path = buf.String()
								}
							}
						}
						exists, sr := ssh.Stat(h, path)
						res = sr
						if task.Register != "" {
							hostVars[task.Register] = map[string]interface{}{"stat": map[string]interface{}{"exists": exists}}
						}
					} else if task.Fail != nil { // fail module
						msg := task.Fail.Msg
						if len(hostVars) > 0 {
							tmpl, err := template.New("fail").Parse(task.Fail.Msg)
							if err == nil {
								var buf bytes.Buffer
								if err := tmpl.Execute(&buf, hostVars); err == nil {
									msg = buf.String()
								}
							}
						}
						res = ssh.CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: msg}
						mu.Lock()
						results = append(results, res)
						mu.Unlock()
						break
					} else if task.Debug != nil { // debug module
						msg := task.Debug.Msg
						if len(hostVars) > 0 {
							tmpl, err := template.New("debug").Parse(task.Debug.Msg)
							if err == nil {
								var buf bytes.Buffer
								if err := tmpl.Execute(&buf, hostVars); err == nil {
									msg = buf.String()
								}
							}
						}
						res = ssh.CommandResult{Host: h.Name, ReturnMsg: "DEBUG", ReturnCode: 0, Output: msg}
					} else {
						res = ssh.CommandResult{
							Host:       h.Name,
							ReturnMsg:  "FAILED",
							ReturnCode: 1,
							Output:     fmt.Sprintf("Unsupported task type in '%s'", task.Name),
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
