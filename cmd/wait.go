package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// waitCmd represents the wait command
var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Waits for completion of a ctzz-managed task",
	Long:  `useful for pausing workflows that depend on task results`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wait called")
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
}
