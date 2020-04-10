package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys all cloud-task-zip-zap infrastructure",
	Long: `destroy all ctzz cloud components`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy called")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
