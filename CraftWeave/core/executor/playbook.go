package executor

import (
	"fmt"
	"sync"

	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

// Executor executes playbooks with configurable behaviour.
type Executor struct {
	AggregateOutput bool
	CheckMode       bool
	MaxWorkers      int
	Logger          LogCollector
}

// New creates a new Executor.
func New(aggregate, check bool) *Executor {
	return &Executor{AggregateOutput: aggregate, CheckMode: check, MaxWorkers: 5}
}

// SetLogger configures a log collector for execution results.
func (e *Executor) SetLogger(l LogCollector) { e.Logger = l }

// Execute processes and runs the given playbook.
func (e *Executor) Execute(playbook []parser.Play, inventoryPath string) {
	for i := range playbook {
		play := &playbook[i]
		if play.Vars == nil {
			play.Vars = make(map[string]string)
		}

		copyVars := func(src map[string]string) map[string]string {
			dst := make(map[string]string, len(src))
			for k, v := range src {
				dst[k] = v
			}
			return dst
		}

		fmt.Printf("\n🎯 Play: %s (hosts: %s)\n", play.Name, play.Hosts)

		hosts, err := inventory.Parse(inventoryPath, play.Hosts)
		if err != nil {
			fmt.Printf("❌ Failed to resolve hosts: %v\n", err)
			continue
		}

		hostVars := make(map[string]map[string]string, len(hosts))
		for _, h := range hosts {
			hostVars[h.Name] = copyVars(play.Vars)
		}

		for _, task := range play.Tasks {
			fmt.Printf("\nTASK [%s] ********************************************************\n", task.Name)

			var results []ssh.CommandResult
			var mu sync.Mutex
			var wg sync.WaitGroup
			sem := make(chan struct{}, e.MaxWorkers)

			for _, host := range hosts {
				vars := hostVars[host.Name]
				wg.Add(1)
				go func(h inventory.Host, vars map[string]string) {
					defer wg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					if !EvaluateWhen(task.When, vars) {
						return
					}

					if e.CheckMode {
						res := ssh.CommandResult{Host: h.Name, ReturnMsg: "SKIPPED", ReturnCode: 0, Output: fmt.Sprintf("dry-run: %s", task.Name)}
						mu.Lock()
						results = append(results, res)
						mu.Unlock()
						if e.Logger != nil {
							e.Logger.Collect(res)
						}
						return
					}

					res := ExecuteTask(task, h, vars)
					mu.Lock()
					results = append(results, res)
					mu.Unlock()
					if e.Logger != nil {
						e.Logger.Collect(res)
					}
				}(host, vars)
			}
			wg.Wait()

			if e.AggregateOutput {
				ssh.AggregatedPrint(results)
			} else {
				for _, r := range results {
					fmt.Printf("%s | %s | rc=%d >>\n%s\n", r.Host, r.ReturnMsg, r.ReturnCode, r.Output)
				}
			}
		}
	}
}
