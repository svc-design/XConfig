// core/executor/playbook.go
package executor

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"

	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

// 全局控制输出样式和 dry-run 模式
var AggregateOutput bool
var CheckMode bool

// ExecutePlaybook 解析并执行整个 playbook
func ExecutePlaybook(playbook []parser.Play, inventoryPath string, extraVars map[string]string) {
	for _, play := range playbook {
		fmt.Printf("\n🎯 Play: %s (hosts: %s)\n", play.Name, play.Hosts)

		hosts, err := inventory.Parse(inventoryPath, play.Hosts)
		if err != nil {
			fmt.Printf("❌ Failed to resolve hosts: %v\n", err)
			continue
		}

		var results []ssh.CommandResult
		var mu sync.Mutex
		var wg sync.WaitGroup

		// merge play vars with extra vars (extra vars override)
		mergedVars := make(map[string]string)
		for k, v := range play.Vars {
			mergedVars[k] = v
		}
		for k, v := range extraVars {
			mergedVars[k] = v
		}

		for _, host := range hosts {
			for _, task := range play.Tasks {
				task := task // 关闭闭包引用
				wg.Add(1)

				go func(h inventory.Host) {
					defer wg.Done()

					if CheckMode {
						fmt.Printf("%s | SKIPPED | dry-run: %s\n", h.Name, task.Name)
						return
					}

					var res ssh.CommandResult
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
