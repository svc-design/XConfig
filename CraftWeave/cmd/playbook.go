// cmd/playbook.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"craftweave/core/executor"
	"craftweave/core/parser"
)

var inventoryPath string

var playbookCmd = &cobra.Command{
	Use:   "playbook [file]",
	Short: "Run a CraftWeave playbook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		fmt.Printf("📜 Executing playbook: %s\n", file)

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("❌ Failed to read playbook: %v\n", err)
			os.Exit(1)
		}

		var plays []parser.Play
		if err := yaml.Unmarshal(data, &plays); err != nil {
			fmt.Printf("❌ Failed to parse YAML: %v\n", err)
			os.Exit(1)
		}

		executor.AggregateOutput = AggregateOutput
		executor.CheckMode = CheckMode

		// parse -e/--extra-vars into map and pass to executor
		mergedVars := make(map[string]string)
		for k, v := range ExtraVars {
			mergedVars[k] = v
		}
		executor.ExecutePlaybook(plays, inventoryPath, mergedVars)
	},
}

func init() {
	playbookCmd.Flags().StringVarP(&inventoryPath, "inventory", "i", "hosts.yaml", "Inventory file")
	playbookCmd.Flags().BoolVarP(&AggregateOutput, "aggregate", "A", false, "Aggregate output from identical results")
	playbookCmd.Flags().BoolVarP(&CheckMode, "check", "C", false, "Dry-run mode")
	playbookCmd.Flags().StringToStringVarP(&ExtraVars, "extra-vars", "e", nil, "Extra variables in key=value format")
	rootCmd.AddCommand(playbookCmd)
}
