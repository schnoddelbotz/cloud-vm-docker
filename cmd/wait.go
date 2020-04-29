package cmd

import (
	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
	"github.com/spf13/cobra"
)

// waitCmd represents the wait command
var waitCmd = &cobra.Command{
	Use:          "wait",
	Short:        "Waits for completion of a cloud-vm-docker-managed task to complete",
	Long:         `useful for pausing workflows that depend on task results; same as --wait on task-vm or not using -d on run`,
	Args:         cobra.RangeArgs(1, 1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		vmID := args[0]
		g := settings.ViperToRuntimeSettings(true)
		_, err := cloud.WaitForTaskDone(g.ProjectID, vmID)
		return err
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
}
