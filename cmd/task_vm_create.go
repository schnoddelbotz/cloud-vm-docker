package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a ComputeEngine VM instance",
	// Does the same like run, but circumvents pubsub; creates DataStore entry and spins up VM
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := viper.GetString(settings.FlagProject)
		zone := viper.GetString(settings.FlagZone)
		image := args[0]
		var command []string
		if len(args) > 1 {
			command = args[1:]
		}
		taskArguments := cloud.NewCloudTaskArgsFromArgs(image,
			command,
			viper.GetString(settings.FlagEntryPoint), // FIXME!!! UNUSED!!!
			viper.GetString(settings.FlagVMType))
		sshKeys := cloud.GetUserSSHPublicKeys(viper.GetString(settings.FlagSSHPublicKey),
			viper.GetBool(settings.FlagNoSSH))

		log.Printf("Writing task to DataStore: %+v", taskArguments)
		task := cloud.StoreNewTask(projectID, *taskArguments)

		createOp, err := cloud.CreateVM(projectID, zone, task, sshKeys)
		if err != nil {
			return fmt.Errorf("ERROR running TaskArguments: %v", err)
		}
		log.Println("VM creation requested successfully")
		cloud.WaitForOperation(projectID, zone, createOp.Name)
		// TODO: Now print('ssh ctzz@IP')
		// https://cloud.google.com/compute/docs/instances/view-ip-address
		// log.Printf("Use ssh ctzz@%s to connect", GetInstanceIP(Name))
		return nil
	},
}

func init() {
	// this should support most flags like 'run'
	flags := createCmd.Flags()
	flags.SetInterspersed(false)

	flags.StringP(settings.FlagVMType, "v", "n1-standard-1", "VM machine type")
	flags.StringP(settings.FlagSSHPublicKey, "s", "", "SSH public key to put on VM (default ~/.ssh/*.pub)")
	flags.BoolP(settings.FlagNoSSH, "n", false, "disable SSH public key install [notyet]")

	viper.BindPFlag(settings.FlagVMType, createCmd.Flags().Lookup(settings.FlagVMType))
	viper.BindPFlag(settings.FlagSSHPublicKey, createCmd.Flags().Lookup(settings.FlagSSHPublicKey))
	viper.BindPFlag(settings.FlagNoSSH, createCmd.Flags().Lookup(settings.FlagNoSSH))

	// add wait flag to wait for spin-up?

	vmCmd.AddCommand(createCmd)
}
