package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows status of one or more ctzz-managed tasks",
	Long:  `fetches status from FireStore, allows filtering state/progress/..., sorting`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called, list all VMs + power state, join tasks?")
	},
}

func init() {
	vmCmd.AddCommand(statusCmd)
}
