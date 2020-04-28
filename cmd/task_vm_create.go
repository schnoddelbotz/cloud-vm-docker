package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a ComputeEngine VM instance",
	// Does the same like run, but circumvents http cfn; creates FireStore entry and spins up GCE VM
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		g := settings.EnvironmentToGoogleSettings(true)
		//e := handlers.NewEnvironment(g, false, true, true)

		image := args[0]
		var command []string
		if len(args) > 1 {
			command = args[1:]
		}
		taskArguments := cloud.NewTaskArgumentsFromArgs(image, command,
			viper.GetString(settings.FlagEntryPoint), // FIXME!!! UNUSED!!!
			g.VMType)

		log.Printf("Writing task to FireStore: %+v", taskArguments)
		task := cloud.StoreNewTask(g.ProjectID, *taskArguments)
		task.SSHPubKeys = cloud.GetUserSSHPublicKeys(g.SSHPublicKey, g.DisableSSH)

		createOp, err := cloud.CreateVM(g, task)
		if err != nil {
			return fmt.Errorf("ERROR running TaskArguments: %v", err)
		}

		log.Printf("VM logs: %s", cloud.GetLogLinkForVM(g.ProjectID, createOp.TargetId))
		log.Println("VM creation requested successfully, waiting for create completion")
		cloud.WaitForOperation(g.ProjectID, g.Zone, createOp.Name)

		err = cloud.SetTaskInstanceId(g.ProjectID, task.VMID, createOp.TargetId)
		if err != nil {
			log.Printf("ARGH!!! Could not update instanceID in FireStore: %s", err)
		}

		if viper.GetBool(settings.FlagWait) {
			nt, err := cloud.WaitForTaskDone(g.ProjectID, task.VMID)
			if err != nil {
				return err
			}

			if viper.GetBool(settings.FlagPrintLogs) {
				log.Printf("WAITING another 15 seconds for logs to appear in StackDriver...")
				time.Sleep(15 * time.Second)
				log.Printf("Logs from StackDriver:")
				cloud.PrintLogEntries(g.ProjectID, nt.InstanceID, nt.DockerContainerId, nt.CreatedAt)
			}
		}

		// TODO: Now print('ssh cloud-vm-docker@IP')
		// https://cloud.google.com/compute/docs/instances/view-ip-address
		// log.Printf("Use ssh cloud-vm-docker@%s to connect", GetInstanceIP(Name))
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
	flags.BoolP(settings.FlagWait, "w", false, "wait until command completes / VM shuts down")
	flags.BoolP(settings.FlagPrintLogs, "P", false, "print (last 512 lines of) container logs (requires --wait)")

	viper.BindPFlag(settings.FlagVMType, createCmd.Flags().Lookup(settings.FlagVMType))
	viper.BindPFlag(settings.FlagSSHPublicKey, createCmd.Flags().Lookup(settings.FlagSSHPublicKey))
	viper.BindPFlag(settings.FlagNoSSH, createCmd.Flags().Lookup(settings.FlagNoSSH))
	viper.BindPFlag(settings.FlagWait, createCmd.Flags().Lookup(settings.FlagWait))
	viper.BindPFlag(settings.FlagPrintLogs, createCmd.Flags().Lookup(settings.FlagPrintLogs))

	// add wait flag to wait for spin-up?

	vmCmd.AddCommand(createCmd)
}
