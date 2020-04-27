package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// AppVersion is set at compile time via make / ldflags
var AppVersion = "0.0.x-dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloud-vm-docker",
	Short: "A brief description of your application",
	Long: `cloud-vm-docker runs Docker images on Google Compute Engine VMs
See https://github.com/schnoddelbotz/cloud-vm-docker`,
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
	rootCmd.PersistentFlags().StringP(settings.FlagZone, "z", "europe-west1-b", "google cloud zone for VMs")
	rootCmd.PersistentFlags().StringP(settings.FlagRegion, "r", "europe-west1", "google cloud region for CloudFunctions")
	rootCmd.PersistentFlags().BoolP(settings.FlagVerbose, "V", false, "enable verbose/debug output of cloud-vm-docker")
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose operations")

	/// DUP!!! lives in settings.go, too FIXME?
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("CVD")
	///
	viper.BindPFlag(settings.FlagProject, rootCmd.PersistentFlags().Lookup(settings.FlagProject))
	viper.BindPFlag(settings.FlagRegion, rootCmd.PersistentFlags().Lookup(settings.FlagRegion))
	viper.BindPFlag(settings.FlagZone, rootCmd.PersistentFlags().Lookup(settings.FlagZone))
	viper.BindPFlag(settings.FlagVerbose, rootCmd.PersistentFlags().Lookup(settings.FlagVerbose))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
