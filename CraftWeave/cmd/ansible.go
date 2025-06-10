package cmd

import (
	"fmt"
	"strings"
	"sync"

	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/modules"
	"craftweave/internal/ssh"
	"github.com/spf13/cobra"
)

var module, args string

var ansibleCmd = &cobra.Command{
	Use:   "ansible [target]",
	Short: "Run ad-hoc tasks on target hosts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, targets []string) {
		hosts, err := inventory.Parse(InventoryPath, targets[0])
		if err != nil {
			fmt.Println("Failed to parse inventory:", err)
			return
		}

		var results []ssh.CommandResult
		var mu sync.Mutex
		var wg sync.WaitGroup
		sem := make(chan struct{}, MaxWorkers)

		for _, h := range hosts {
			wg.Add(1)
			go func(h inventory.Host) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				if CheckMode {
					fmt.Printf("%s | SKIPPED\n", h.Name)
					return
				}

				t := parser.Task{}
				switch module {
				case "shell":
					t.Shell = args
				case "script":
					t.Script = args
				case "template":
					parts := strings.SplitN(args, ":", 2)
					if len(parts) == 2 {
						t.Template = &parser.Template{Src: parts[0], Dest: parts[1]}
					}
				}

				if hdl, ok := modules.GetHandler(t.Type()); ok {
					res := hdl(modules.Context{Host: h}, t)
					mu.Lock()
					results = append(results, res)
					mu.Unlock()
				} else {
					mu.Lock()
					results = append(results, ssh.CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("Module '%s' is not supported.", module)})
					mu.Unlock()
				}
			}(h)
		}
		wg.Wait()

		if AggregateOutput {
			ssh.AggregatedPrint(results)
		} else {
			for _, r := range results {
				fmt.Printf("%s | %s | rc=%d >>\n%s\n", r.Host, r.ReturnMsg, r.ReturnCode, r.Output)
			}
		}
	},
}

func init() {
	ansibleCmd.Flags().StringVarP(&InventoryPath, "inventory", "i", "hosts.yaml", "Inventory file")
	ansibleCmd.Flags().StringVarP(&module, "module", "m", "shell", "Module to execute")
	ansibleCmd.Flags().StringVarP(&args, "args", "a", "", "Arguments for the module")
	ansibleCmd.Flags().IntVarP(&MaxWorkers, "forks", "f", 5, "Max parallel tasks")
	ansibleCmd.Flags().BoolVarP(&CheckMode, "check", "C", false, "Check mode (dry-run)")
	ansibleCmd.Flags().BoolVarP(&AggregateOutput, "aggregate", "A", false, "Aggregate identical output")
	rootCmd.AddCommand(ansibleCmd)
}
