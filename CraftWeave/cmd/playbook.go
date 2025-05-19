package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var playbookCmd = &cobra.Command{
	Use:   "playbook [file]",
	Short: "Run a CraftWeave playbook",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		fmt.Printf("[CraftWeave] Running playbook: %s\n", file)
		// TODO: 加载并执行 playbook
	},
}
