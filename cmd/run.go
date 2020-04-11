package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a dockerized command to be executed on a ComputeEngine VM",
	Long:  `run dockerized command, pass args via json (from cmd line or file)`,
	Run: func(cmd *cobra.Command, args []string) {
		task := cloud.NewCloudTaskFromArgs(viper.GetString(settings.FlagImage),
			viper.GetString(settings.FlagCommand),
			viper.GetString(settings.FlagEntryPoint),
			viper.GetString(settings.FlagVMType))
		err := cloud.PubSubPushTask(task, viper.GetString(settings.FlagProject), settings.TopicNameTaskQueue)
		if err != nil {
			fmt.Printf("ERROR publishing task: %v\n", err)
		}
		fmt.Printf("SUCCESS\n")
	},
}

func init() {
	runCmd.Flags().StringP(settings.FlagImage, "i", "busybox", "image to run on VM")
	runCmd.Flags().StringP(settings.FlagCommand, "c", "", "command to run in container")
	//runCmd.Flags().StringP(settings.FlagArgs, "a", "{}", "JSON args to pass to container app") -- complicates, not needed, command can do all?
	runCmd.Flags().StringP(settings.FlagArgsFile, "f", "", "like --args, but read from given file")
	// todo: have bool flag --via-cfn -- as this client runs directly into pubsub by default.
	//       "legacy" scripts may want to spawn VMs just via CFN/HTTP+Token, to avoid need for this binary (+svc_acc)

	viper.BindPFlag(settings.FlagImage, runCmd.Flags().Lookup(settings.FlagImage))
	viper.BindPFlag(settings.FlagCommand, runCmd.Flags().Lookup(settings.FlagCommand))
	//viper.BindPFlag(settings.FlagArgs, runCmd.Flags().Lookup(settings.FlagArgs))
	viper.BindPFlag(settings.FlagArgsFile, runCmd.Flags().Lookup(settings.FlagArgsFile))

	rootCmd.AddCommand(runCmd)
}
