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
		// https://github.com/spf13/viper/issues/233#issuecomment-479336184
		viper.BindPFlag(settings.FlagWait, cmd.Flags().Lookup(settings.FlagWait))
		viper.BindPFlag(settings.FlagPrintLogs, cmd.Flags().Lookup(settings.FlagPrintLogs))

		// shared with runCmd:
		viper.BindPFlag(settings.FlagNoSSH, cmd.Flags().Lookup(settings.FlagNoSSH))
		viper.BindPFlag(settings.FlagSSHPublicKey, cmd.Flags().Lookup(settings.FlagSSHPublicKey))
		viper.BindPFlag(settings.FlagVMType, cmd.Flags().Lookup(settings.FlagVMType))
		viper.BindPFlag(settings.FlagSubnet, cmd.Flags().Lookup(settings.FlagSubnet))
		viper.BindPFlag(settings.FlagTags, cmd.Flags().Lookup(settings.FlagTags))

		g := settings.EnvironmentToGoogleSettings(true)
		//e := handlers.NewEnvironment(g, false, true, true)

		image := args[0]
		var command []string
		if len(args) > 1 {
			command = args[1:]
		}
		taskArguments := cloud.NewTaskArgumentsFromArgs(image, command,
			viper.GetString(settings.FlagEntryPoint), // FIXME!!! UNUSED!!!
			g.VMType, viper.GetString(settings.FlagSubnet), viper.GetString(settings.FlagTags))

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
			log.Printf("Docker container logs: %s", cloud.GetLogLinkForContainer(g.ProjectID, createOp.TargetId, nt.DockerContainerId))
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

	flags.BoolP(settings.FlagWait, "w", false, "wait until command completes / VM shuts down")
	flags.BoolP(settings.FlagPrintLogs, "P", false, "print (last 512 lines of) container logs (requires --wait)")

	// shared with runCmd: -- todo: better way?
	flags.BoolP(settings.FlagNoSSH, "n", false, "disable SSH public key install [notyet]")
	flags.StringP(settings.FlagSSHPublicKey, "s", "", "SSH public key to put on VM (default ~/.ssh/*.pub)")
	flags.StringP(settings.FlagVMType, "v", "n1-standard-1", "VM machine type")
	flags.StringP(settings.FlagSubnet, "S", "", "optional non-default subnet for VM")
	flags.StringP(settings.FlagTags, "T", "", "VM tags (comma-separated list of tags")

	vmCmd.AddCommand(createCmd)
}
