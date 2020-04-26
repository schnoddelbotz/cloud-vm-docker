package cmd

import (
	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

// logsCmd represents the logs command, yay!
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Downloads or prints VM logs",
	Long:  `download and print VM logs by task UUID`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vmID := args[0]
		g := settings.EnvironmentToGoogleSettings(true)

		task, err := cloud.GetTask(g.ProjectID, vmID)
		if err != nil {
			log.Fatalf("Unable to get task data for vm_id %s: %s", vmID, err)
		}
		instanceID, err := strconv.ParseUint(task.InstanceID, 10, 64)
		if err != nil {
			log.Fatalf("Oops, unable to convert instanceID %s to uint64: %s", task.InstanceID, err)
		}
		if viper.GetBool(settings.FlagVerbose) {
			log.Printf("VM Logs    : %s", cloud.GetLogLinkForVM(g.ProjectID, instanceID))
			log.Printf("Docker Logs: %s", cloud.GetLogLinkForContainer(g.ProjectID, instanceID, task.DockerContainerId))
		}

		cloud.PrintLogEntries(g.ProjectID, task.InstanceID, task.DockerContainerId, task.CreatedAt)
	},
}

func init() {
	vmCmd.AddCommand(logsCmd)
	// based on task UUID, looks up VM stackdriver logs and streams to stdout or file
	// --follow ? -- until proc exits?
	// --to-file
	// --log-fmt <google-fmt>
	// --custom-quer ...
	// --sort asc|desc
	// --limit
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
