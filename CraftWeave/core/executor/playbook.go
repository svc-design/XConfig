package executor

import (
	"fmt"

	"craftweave/core/output"
	"craftweave/core/parser"
	"craftweave/internal/inventory"
)

// Global options
var AggregateOutput bool
var CheckMode bool

// ExecutePlaybook parses and runs a playbook using registered task modules
func ExecutePlaybook(playbook []parser.Play, inventoryPath string, concurrency int) {
	for _, play := range playbook {
		fmt.Printf("\n🎯 Play: %s (hosts: %s)\n", play.Name, play.Hosts)

		hosts, err := inventory.Parse(inventoryPath, play.Hosts)
		if err != nil {
			fmt.Printf("❌ Failed to resolve hosts: %v\n", err)
			continue
		}

		var collector output.Collector
		if AggregateOutput {
			collector = &output.AggregateCollector{}
		} else {
			collector = output.StdoutCollector{}
		}

		pool := NewPool(concurrency)
		for _, host := range hosts {
			for _, task := range play.Tasks {
				h := host
				t := task
				pool.Go(func() {
					if CheckMode {
						fmt.Printf("%s | SKIPPED | dry-run: %s\n", h.Name, t.Name)
						return
					}
					res := ExecuteTask(h, t, play.Vars)
					collector.Collect(res)
				})
			}
		}
		pool.Wait()
		collector.Flush()
	}
}
