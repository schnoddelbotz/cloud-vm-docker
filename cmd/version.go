package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Prints the cloud-vm-docker version",
		Run: func(cmd *cobra.Command, args []string) {
			println(AppVersion)
		},
	})
}
