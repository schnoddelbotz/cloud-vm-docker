package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// vmCmd represents the vm command
var vmCmd = &cobra.Command{
	Use:   "task-vm",
	Short: "Low-level ctzz VM management (create, kill, status ...)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vm called")
	},
}

func init() {
	rootCmd.AddCommand(vmCmd)
}
