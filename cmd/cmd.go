package cmd

import (
	"os"

	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// AppVersion is set at compile time via make / ldflags
var AppVersion = "0.0.x-dev"

// EnvironmentToSettings translates environment variables into a Settings struct.
// For CLI gotsfn, this is done by cobra/viper.
func EnvironmentToSettings() settings.Settings {
	s := settings.Settings{
		Provider: os.Getenv(settings.ActionKill),
		//GCS: settings.GCSSettings{
		//	// projectID is set from the GCP_PROJECT environment variable, which is
		//	// automatically set by the Cloud Functions runtime.
		//	//ProjectID:  os.Getenv(settings.EnvKeyGCPProjectID),
		//	BucketName: os.Getenv(settings.EnvKeyGCSBucketName),
		//	//BucketRoot: os.Getenv(settings.EnvKeyGCPBucketRoot), // not yet supported
		//},
		//S3: settings.S3Settings{
		//	AccountID: os.Getenv(settings.EnvKeyAWSAccountID),
		//},
		//Filesystem: settings.FilesystemSettings{
		//	Directory: os.Getenv(settings.EnvKeyDirectory),
		//},
	}
	return s
}
