package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys all cloud-vm-docker cloud infrastructure",
	Long:  `destroy all cloud-vm-docker cloud components`,
	Run: func(cmd *cobra.Command, args []string) {
		cloud.Destroy(viper.GetString(settings.FlagProject))
	},
}

func init() {
	setupCmd.AddCommand(destroyCmd)
}
