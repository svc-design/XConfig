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

		for _, host := range hosts {
			for _, task := range allTasks {
				task := task // 关闭闭包引用
				wg.Add(1)

				go func(h inventory.Host) {
					defer wg.Done()

					if CheckMode {
						fmt.Printf("%s | SKIPPED | dry-run: %s\n", h.Name, task.Name)
						return
					}

					var res ssh.CommandResult
					mergedVars := make(map[string]string)
					for k, v := range hostFacts[h.Name] {
						mergedVars[k] = v
					}
					for k, v := range playVars {
						mergedVars[k] = v
					}

					if task.Shell != "" {
						rendered := task.Shell
						if len(mergedVars) > 0 {
							renderedTmpl, err := template.New("shell").Parse(task.Shell)
							if err == nil {
								var buf bytes.Buffer
								if err := renderedTmpl.Execute(&buf, mergedVars); err == nil {
									rendered = buf.String()
								}
							}
						}
						res = ssh.RunShellCommand(h, rendered)
					} else if task.Script != "" {
						res = ssh.RunRemoteScript(h, task.Script)
					} else if task.Template != nil {
						res = ssh.RenderTemplate(h, task.Template.Src, task.Template.Dest, mergedVars)
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
				}(host)
			}
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
