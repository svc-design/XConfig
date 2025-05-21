package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

var inventoryPath, module, args string
var check bool

var ansibleCmd = &cobra.Command{
	Use:   "ansible [target]",
	Short: "Run ad-hoc tasks on target hosts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, targets []string) {
		hosts, err := inventory.Parse(inventoryPath, targets[0])
		if err != nil {
			fmt.Println("Failed to parse inventory:", err)
			return
		}

		if module != "shell" {
			fmt.Printf("Only 'shell' module is implemented. '%s' is not supported.\n", module)
			return
		}

		var wg sync.WaitGroup
		for _, h := range hosts {
			wg.Add(1)
			go func(h inventory.Host) {
				defer wg.Done()
				if check {
					fmt.Printf("%s | SKIPPED\n", h.Name)
					return
				}
				ssh.RunShellCommand(h, args)
			}(h)
		}
		wg.Wait()
	},
}

func init() {
	ansibleCmd.Flags().StringVarP(&inventoryPath, "inventory", "i", "hosts.yaml", "Inventory file")
	ansibleCmd.Flags().StringVarP(&module, "module", "m", "shell", "Module to execute")
	ansibleCmd.Flags().StringVarP(&args, "args", "a", "", "Arguments for the module")
	ansibleCmd.Flags().BoolVarP(&check, "check", "C", false, "Check mode (dry-run)")
	rootCmd.AddCommand(ansibleCmd)
}
