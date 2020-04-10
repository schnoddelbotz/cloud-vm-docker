package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys all ctzz cloud infrastructure",
	Long:  `destroy all ctzz cloud components`,
	Run: func(cmd *cobra.Command, args []string) {
		cloud.Destroy(viper.GetString(settings.FlagProject))
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
