package cmd

import (
	"fmt"
	"sync"

	"craftweave/internal/inventory"
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

		for _, h := range hosts {
			wg.Add(1)
			go func(h inventory.Host) {
				defer wg.Done()
				if CheckMode {
					fmt.Printf("%s | SKIPPED\n", h.Name)
					return
				}

				var res ssh.CommandResult
				switch module {
				case "shell":
					res = ssh.RunShellCommand(h, args)
				case "command":
					res = ssh.RunCommand(h, args)
				case "script":
					res = ssh.RunRemoteScript(h, args)
				default:
					res = ssh.CommandResult{
						Host:       h.Name,
						ReturnMsg:  "FAILED",
						ReturnCode: 1,
						Output:     fmt.Sprintf("Module '%s' is not supported.", module),
					}
				}

				mu.Lock()
				results = append(results, res)
				mu.Unlock()
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
	ansibleCmd.Flags().BoolVarP(&CheckMode, "check", "C", false, "Check mode (dry-run)")
	ansibleCmd.Flags().BoolVarP(&AggregateOutput, "aggregate", "A", false, "Aggregate identical output")
	rootCmd.AddCommand(ansibleCmd)
}
