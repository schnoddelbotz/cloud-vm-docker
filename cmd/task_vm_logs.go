package cmd

import (
	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
	"github.com/spf13/cobra"
	"log"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Downloads or prints VM logs",
	Long:  `download and print VM logs by task UUID`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// todo! vm name (aka vmID, string) vs. GCE instance id (uint64)
		var vmID uint64 = 3517187272
		g := settings.EnvironmentToGoogleSettings(true)
		log.Printf("DEBUG:\n VM LogLink: %s\n ContainerLogLink: %s",
			cloud.GetLogLinkForVM(g.ProjectID, vmID),
			cloud.GetLogLinkForContainer(g.ProjectID, vmID, "aaa-blah-test"))
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
