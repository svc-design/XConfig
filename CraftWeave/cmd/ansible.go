package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var inventory string
var module string
var args string

var ansibleCmd = &cobra.Command{
	Use:   "ansible [target]",
	Short: "Run ad-hoc tasks on target hosts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, target []string) {
		host := target[0] // e.g. "all"
		fmt.Printf("[CraftWeave] Running module '%s' on host '%s'\n", module, host)

		if module == "ping" {
			// Just a local ping test to simulate remote call
			cmd := exec.Command("ping", "-c", "1", "127.0.0.1")
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(string(out))
			}
		} else {
			fmt.Printf("Module '%s' is not implemented yet.\n", module)
		}
	},
}

func init() {
	ansibleCmd.Flags().StringVarP(&inventory, "inventory", "i", "hosts.yaml", "Inventory file")
	ansibleCmd.Flags().StringVarP(&module, "module", "m", "ping", "Module to execute")
	ansibleCmd.Flags().StringVarP(&args, "args", "a", "", "Arguments for the module")
}
