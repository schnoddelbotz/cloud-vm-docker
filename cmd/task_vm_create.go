package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a ComputeEngine VM instance",
	// Does the same like run, but circumvents pubsub; creates DataStore entry and spins up VM
	//Run: func(cmd *cobra.Command, args []string) {
	//	// FIXME: this should pass a task
	//	taskArgs := cloud.NewCloudTaskArgsFromArgs(viper.GetString(settings.))
	//	cloud.CreateVM(viper.GetString(settings.FlagProject),
	//		viper.GetString(settings.FlagZone),
	//		viper.GetString(settings.FlagVMType),
	//		"fixme-p")
	//},
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := args[0]
		var command []string
		if len(args) > 1 {
			command = args[1:]
		}
		taskArguments := cloud.NewCloudTaskArgsFromArgs(image,
			command,
			viper.GetString(settings.FlagEntryPoint), // FIXME!!! UNUSED!!!
			viper.GetString(settings.FlagVMType))
		log.Printf("Writing task to DataStore: %+v", taskArguments)
		task := cloud.StoreNewTask(viper.GetString(settings.FlagProject), *taskArguments)
		log.Printf("Success writing Datastore, now creating VM..")
		createOp, err := cloud.CreateVM(viper.GetString(settings.FlagProject),
			viper.GetString(settings.FlagZone),
			task)
		if err != nil {
			return fmt.Errorf("ERROR running TaskArguments: %v", err)
		}
		log.Println("VM Created successfully!")
		log.Printf("  Operation   : %s", createOp)
		log.Printf("  SSH         : to-do ...")
		return nil
	},
}

func init() {
	flags := createCmd.Flags()
	flags.SetInterspersed(false)
	flags.StringP(settings.FlagVMType, "v", "n1-standard-1", "VM machine type")
	viper.BindPFlag(settings.FlagVMType, createCmd.Flags().Lookup(settings.FlagVMType))
	// add wait flag to wait for spin-up?
	vmCmd.AddCommand(createCmd)
}
