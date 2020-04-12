package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// vmCmd represents the vm command
var vmCmd = &cobra.Command{
	Use:   "task-vm",
	Short: "Low-level ctzz VM management (create, kill, status ...)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vm called")
	},
}

func init() {
	vmCmd.Flags().StringP(settings.FlagVMType, "v", "n1-standard-1", "VM machine type")
	viper.BindPFlag(settings.FlagVMType, runCmd.Flags().Lookup(settings.FlagVMType))
	rootCmd.AddCommand(vmCmd)
}
