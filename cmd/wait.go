package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// waitCmd represents the wait command
var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Waits for completion of a cloud-vm-docker-managed task to complete",
	Long:  `useful for pausing workflows that depend on task results`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wait called ... should wait for one or more VMs to report COMPLETED_*")
		// verbose should tell which tasks found in which state...
		// this is a pure datastore op... but could be http-cfn'ed as well
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
}
