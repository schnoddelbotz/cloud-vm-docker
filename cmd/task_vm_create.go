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
	// Does the same like run, but circumvents http cfn; creates DataStore entry and spins up GCE VM
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

		log.Printf("Writing task to DataStore: %+v", taskArguments)
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
			log.Printf("ARGH!!! Could not update instanceID in DataStore: %s", err)
		}

		status := "submitted"
		if viper.GetBool(settings.FlagWait) {
			// FIXME extract ... and avoid polling if possible.
			var t cloud.Task
			log.Printf("Waiting for command completion now...")
			for status != "DONE" {
				// fixme call CFN, not DS directly???!
				t, err = cloud.GetTask(g.ProjectID, task.VMID)
				if err != nil {
					log.Fatalf("Failed to poll DS: %s", err)
				}
				if status != t.Status {
					log.Printf("Status changed from %s -> %s", status, t.Status)
				}
				status = t.Status
				if status != "DONE" {
					time.Sleep(30 * time.Second)
				}
			}
			log.Printf("Container logs: %s", cloud.GetLogLinkForContainer(g.ProjectID, createOp.TargetId, t.DockerContainerId))
			log.Printf("Task exit code: %d", t.DockerExitCode)
			if t.DockerExitCode != 0 {
				return fmt.Errorf("task exited with non-zero return code %d", t.DockerExitCode)
			}

			if viper.GetBool(settings.FlagPrintLogs) {
				log.Printf("WAITING another 10 seconds for logs to appear in StackDriver...")
				time.Sleep(10 * time.Second)
				log.Printf("Logs from StackDriver:")
				cloud.PrintLogEntries(g.ProjectID, t.InstanceID, t.DockerContainerId, t.CreatedAt)
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
