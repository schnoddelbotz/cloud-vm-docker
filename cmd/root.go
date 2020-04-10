package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloud-task-zip-zap",
	Short: "A brief description of your application",
	Long: `cloud-task-zip-zap (ctzz) runs Docker images on Google Compute Engine VMs
See https://github.com/schnoddelbotz/cloud-task-zip-zap`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP(settings.FlagProject, "p", "", "google cloud project to work within")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("CTZZ")
	viper.BindPFlag(settings.FlagProject, rootCmd.PersistentFlags().Lookup(settings.FlagProject))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
