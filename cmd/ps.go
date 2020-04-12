package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Prints VM instances, their containers, and status (like docker ps)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ps called")
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolP("all", "a", false, "print deleted VMs, too")
	// tbd: add ctzz system prune to delete ... stuff.
}
