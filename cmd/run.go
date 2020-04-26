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
		g := settings.EnvironmentToGoogleSettings(false)
		taskArguments := cloud.NewTaskArgumentsFromArgs(image, command,
			viper.GetString(settings.FlagEntryPoint), g.VMType)
		endpoint := cloud.GetEndpoint(viper.GetString(settings.FlagProject), viper.GetString(settings.FlagRegion))
		client := cloud.NewCFNClient(endpoint, viper.GetString(settings.FlagToken))

		taskData, err := client.Run(*taskArguments)
		if err != nil {
			return fmt.Errorf("ERROR running TaskArguments: %v", err)
		}

		if !viper.GetBool(settings.FlagDetached) {
			client.WaitForDoneStatus(taskData.VMID)
		}

		return nil
	},
}

func init() {
	// https://github.com/docker/cli/blob/master/cli/command/container/run.go
	flags := runCmd.Flags()
	flags.SetInterspersed(false)

	flags.BoolP(settings.FlagDetached, "d", false, "Start and directly return container ID (and quit)")
	flags.StringP(settings.FlagToken, "t", "", "CloudVMDocker HTTP CloudFunction access token")
	flags.StringP(settings.FlagVMType, "v", "n1-standard-1", "VM machine type")

	viper.BindPFlag(settings.FlagDetached, runCmd.Flags().Lookup(settings.FlagDetached))
	viper.BindPFlag(settings.FlagToken, runCmd.Flags().Lookup(settings.FlagToken))
	viper.BindPFlag(settings.FlagVMType, runCmd.Flags().Lookup(settings.FlagVMType))

	rootCmd.AddCommand(runCmd)
}
