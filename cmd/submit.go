package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a dockerized command to be executed on a ComputeEngine VM",
	Long:  `run dockerized command, pass args via json (from cmd line or file)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("submit called")
	},
}

func init() {
	submitCmd.Flags().StringP(settings.FlagImage, "i", "busybox", "image to run on VM")
	submitCmd.Flags().StringP(settings.FlagCommand, "c", "", "command to run in container")
	submitCmd.Flags().StringP(settings.FlagArgs, "a", "{}", "JSON args to pass to container app")
	submitCmd.Flags().StringP(settings.FlagArgsFile, "f", "", "like --args, but read from given file")
	// todo: have bool flag --via-cfn -- as this client submits directly into pubsub by default.
	//       "legacy" scripts may want to spawn VMs just via CFN/HTTP+Token, to avoid need for this binary (+svc_acc)

	viper.BindPFlag(settings.FlagImage, submitCmd.Flags().Lookup(settings.FlagImage))
	viper.BindPFlag(settings.FlagCommand, submitCmd.Flags().Lookup(settings.FlagCommand))
	viper.BindPFlag(settings.FlagArgs, submitCmd.Flags().Lookup(settings.FlagArgs))
	viper.BindPFlag(settings.FlagArgsFile, submitCmd.Flags().Lookup(settings.FlagArgsFile))

	rootCmd.AddCommand(submitCmd)
}
