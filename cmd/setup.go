package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up cloud-vm-docker CloudFunction/PubSub environment",
	Long:  `cloud-vm-docker setup wizard`,
	Run: func(cmd *cobra.Command, args []string) {
		cloud.Setup(viper.GetString(settings.FlagProject))
	},
}

func init() {
	// add option do disable http cloudfunctions -- more secure, but no curl et al, only this tool or gcloud invoke
	// DISABLED for now ...
	// rootCmd.AddCommand(setupCmd)
}
