package cmd

import (
	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up ctzz CloudFunction/PubSub environment",
	Long: `cloud-task-zip-zap setup wizard`,
	Run: func(cmd *cobra.Command, args []string) {
		cloud.Setup()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
