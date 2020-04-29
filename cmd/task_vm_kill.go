package cmd

import (
	"fmt"
	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:          "kill",
	Short:        "Kills a cloud-vm-docker-managed VM",
	Long:         `kill calls ComputeEngine API and requests instant VM deletion`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceName := args[0]
		g := settings.ViperToRuntimeSettings(true)
		if err := cloud.DeleteInstanceByName(g, instanceName); err != nil {
			return err
		}
		// FIXME! The public Delete function should also update FireStore!
		fmt.Println(instanceName)
		return nil
	},
}

func init() {
	vmCmd.AddCommand(killCmd)
}
