package cmd

import (
	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Prints VM instances, their containers, and status (like docker ps)",
	Run: func(cmd *cobra.Command, args []string) {
		task := cloud.Task{
			Image:      "busybox2",
			Command:    nil,
			EntryPoint: "ep2",
			VMType:     "vmt2",
		}
		cloud.StoreTask(viper.GetString(settings.FlagProject), task)
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolP("all", "a", false, "print deleted VMs, too")
	// tbd: add ctzz system prune to delete ... stuff.
}
