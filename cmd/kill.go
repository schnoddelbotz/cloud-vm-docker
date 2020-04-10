package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kills a ctzz-managed VM",
	Long:  `kill calls ComputeEngine API and requests instant VM deletion`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kill called")
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
