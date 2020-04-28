package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/api_client"
	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run IMAGE [COMMAND] [ARG...]",
	Short: "run a dockerized command to be executed on a ComputeEngine VM",
	Long: `run dockerized command on ComputeEngine VM
Despite Usage message below, no cloud-vm-docker [flags] are supported after [COMMAND] [ARG...]`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		runBegin := time.Now()
		log.Printf(`cloud-vm-docker version %s starting in "run" mode (using --token to talk to CloudFunction)`, AppVersion)
		// https://github.com/spf13/viper/issues/233#issuecomment-479336184
		viper.BindPFlag(settings.FlagDetached, cmd.Flags().Lookup(settings.FlagDetached))
		viper.BindPFlag(settings.FlagToken, cmd.Flags().Lookup(settings.FlagToken))
		// shared with createCmd:
		viper.BindPFlag(settings.FlagNoSSH, cmd.Flags().Lookup(settings.FlagNoSSH))
		viper.BindPFlag(settings.FlagSSHPublicKey, cmd.Flags().Lookup(settings.FlagSSHPublicKey))
		viper.BindPFlag(settings.FlagVMType, cmd.Flags().Lookup(settings.FlagVMType))
		viper.BindPFlag(settings.FlagSubnet, cmd.Flags().Lookup(settings.FlagSubnet))
		viper.BindPFlag(settings.FlagTags, cmd.Flags().Lookup(settings.FlagTags))

		image := args[0]
		var command []string
		if len(args) > 1 {
			command = args[1:]
		}
		g := settings.EnvironmentToGoogleSettings(false)
		taskArguments := cloud.NewTaskArgumentsFromArgs(image, command,
			viper.GetString(settings.FlagEntryPoint), g.VMType, viper.GetString(settings.FlagSubnet), viper.GetString(settings.FlagTags))
		endpoint := api_client.GetEndpoint(viper.GetString(settings.FlagProject), viper.GetString(settings.FlagRegion))
		client := api_client.NewCFNClient(endpoint, viper.GetString(settings.FlagToken))

		taskData, err := client.Run(*taskArguments)
		if err != nil {
			return fmt.Errorf("ERROR running TaskArguments: %v", err)
		}

		instanceID, err := strconv.ParseUint(taskData.InstanceID, 10, 64)
		if err != nil {
			return fmt.Errorf("oops, unable to convert instanceID %s to uint64: %s", taskData.InstanceID, err)
		}
		log.Printf("VM logs: %s", cloud.GetLogLinkForVM(g.ProjectID, instanceID))

		if !viper.GetBool(settings.FlagDetached) {
			task := client.WaitForDoneStatus(taskData.VMID)
			log.Printf("Docker container logs: %s", cloud.GetLogLinkForContainer(g.ProjectID, instanceID, task.DockerContainerId))
			log.Printf("Task execution took %.0f seconds", time.Now().Sub(runBegin).Seconds())
			log.Printf("Docker container exit code: %d", task.DockerExitCode)
			if task.DockerExitCode != 0 {
				return fmt.Errorf("non-zero exit code from container: %d", task.DockerExitCode)
			}
			return nil
		}

		log.Printf("Task submission took %.0f seconds", time.Now().Sub(runBegin).Seconds())
		println(taskData.VMID)
		return nil
	},
}

func init() {
	// https://github.com/docker/cli/blob/master/cli/command/container/run.go
	flags := runCmd.Flags()
	flags.SetInterspersed(false)

	flags.BoolP(settings.FlagDetached, "d", false, "Start and directly return container ID (and quit)")
	flags.StringP(settings.FlagToken, "t", "", "CloudVMDocker HTTP CloudFunction access token")

	// shared with createCmd /  task-vm create:
	flags.BoolP(settings.FlagNoSSH, "n", false, "disable SSH public key install [notyet]")
	flags.StringP(settings.FlagSSHPublicKey, "s", "", "SSH public key to put on VM (default ~/.ssh/*.pub)")
	flags.StringP(settings.FlagVMType, "v", "n1-standard-1", "VM machine type")
	flags.StringP(settings.FlagSubnet, "S", "", "optional non-default subnet for VM")
	flags.StringP(settings.FlagTags, "T", "", "VM tags (comma-separated list of tags")

	rootCmd.AddCommand(runCmd)
}
