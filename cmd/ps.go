package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Prints VM instances, their containers, and status (like docker ps)",
	Run: func(cmd *cobra.Command, args []string) {
		cloud.ListTasks(viper.GetString(settings.FlagProject))
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolP("all", "a", false, "print deleted VMs, too")
	// tbd: add cloud-vm-docker system prune to delete ... stuff.
	// todo: add columns with cpu usage etc
	// todo: let user select columns to display, like vm-type...
	// todo: ALSO FETCH REAL VM (power) STATE!!!!
}
