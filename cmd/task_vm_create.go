package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a ComputeEngine VM instance",

	Run: func(cmd *cobra.Command, args []string) {
		cloud.CreateVM(viper.GetString(settings.FlagProject),
			viper.GetString(settings.FlagZone),
			viper.GetString(settings.FlagVMType),
			"fixme-p")
	},
}

func init() {
	vmCmd.AddCommand(createCmd)
}
