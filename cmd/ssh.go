package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSHs into given VM instance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssh called -- should looku p ip, exec /usr/bin/ssh user@ip @argv")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
