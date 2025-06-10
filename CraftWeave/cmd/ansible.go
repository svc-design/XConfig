package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"craftweave/core/executor"
	"craftweave/core/output"
	"craftweave/core/parser"
	"craftweave/internal/inventory"
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

		var task parser.Task
		switch module {
		case "shell":
			task.Shell = args
		case "script":
			task.Script = args
		default:
			fmt.Printf("Module '%s' is not supported.\n", module)
			return
		}

		var collector output.Collector
		if AggregateOutput {
			collector = &output.AggregateCollector{}
		} else {
			collector = output.StdoutCollector{}
		}

		pool := executor.NewPool(MaxConcurrency)
		for _, h := range hosts {
			host := h
			pool.Go(func() {
				if CheckMode {
					fmt.Printf("%s | SKIPPED\n", host.Name)
					return
				}
				res := executor.ExecuteTask(host, task, nil)
				collector.Collect(res)
			})
		}
		pool.Wait()
		collector.Flush()
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
