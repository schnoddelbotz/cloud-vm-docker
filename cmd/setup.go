package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up ctzz CloudFunction/PubSub environment",
	Long:  `cloud-task-zip-zap setup wizard`,
	Run: func(cmd *cobra.Command, args []string) {
		cloud.Setup(viper.GetString(settings.FlagProject))
	},
}

func init() {
	// add option do disable http cloudfunctions -- more secure, but no curl et al, only this tool or gcloud invoke
	rootCmd.AddCommand(setupCmd)
}
