package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "craftweave",
	Short: "CraftWeave - 执行与编织任务和架构的现代工具",
	Long:  `CraftWeave 是一个现代化的 DevOps CLI 工具，融合任务执行、架构编排、拓扑建模与插件生态。`,
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(ansibleCmd)
	rootCmd.AddCommand(playbookCmd)
	rootCmd.AddCommand(vaultCmd)
	rootCmd.AddCommand(cmdbCmd)
	rootCmd.AddCommand(pluginCmd)
}

func printBanner() {
	content, err := os.ReadFile("banner.txt")
	if err == nil {
		fmt.Println(string(content))
	}
}

