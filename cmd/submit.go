package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a dockerized command to be executed on a ComputeEngine VM",
	Long: `run dockerized command, pass args via json (from cmd line or file)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("submit called")
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
