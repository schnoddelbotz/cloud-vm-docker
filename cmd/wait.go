package cmd

import (
	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
	"log"

	"github.com/spf13/cobra"
)

// waitCmd represents the wait command
var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Waits for completion of a cloud-vm-docker-managed task to complete",
	Long:  `useful for pausing workflows that depend on task results; same as --wait on task-vm or not using -d on run`,
	Args:  cobra.RangeArgs(1, 1), //  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vmID := args[0]
		g := settings.EnvironmentToGoogleSettings(true)
		err := cloud.WaitForTaskDone(g.ProjectID, vmID)
		if err != nil {
			log.Fatalf("Unable to wait for task of vm_id %s: %s", vmID, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
}
