package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run IMAGE [COMMAND] [ARG...]",
	Short: "run a dockerized command to be executed on a ComputeEngine VM",
	Long: `run dockerized command on ComputeEngine VM
Despite Usage message below, no cloud-vm-docker [flags] are supported after [COMMAND] [ARG...]`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := args[0]
		var command []string
		if len(args) > 1 {
			command = args[1:]
		}
		g := settings.EnvironmentToGoogleSettings()
		taskArguments := cloud.NewTaskArgumentsFromArgs(image, command,
			viper.GetString(settings.FlagEntryPoint), g.VMType)
		client := cloud.NewCFNClient("", "")
		_, err := client.Run(*taskArguments)
		if err != nil {
			return fmt.Errorf("ERROR running TaskArguments: %v", err)
		}
		return nil
	},
}

func init() {
	// https://github.com/docker/cli/blob/master/cli/command/container/run.go
	flags := runCmd.Flags()
	flags.SetInterspersed(false)
	flags.StringP(settings.FlagDetached, "d", "detach", "Run container in background and print container ID")
	flags.StringP(settings.FlagVMType, "v", "n1-standard-1", "VM machine type")
	flags.BoolP(settings.FlagWait, "w", false, "wait until command completes / VM shuts down")

	viper.BindPFlag(settings.FlagWait, runCmd.Flags().Lookup(settings.FlagWait))
	viper.BindPFlag(settings.FlagDetached, runCmd.Flags().Lookup(settings.FlagDetached))
	viper.BindPFlag(settings.FlagVMType, runCmd.Flags().Lookup(settings.FlagVMType))

	rootCmd.AddCommand(runCmd)
}
