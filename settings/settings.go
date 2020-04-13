package settings

import (
	"strings"

	"github.com/spf13/viper"
)

// GoogleSettings define anything Google related (project, service account, ...)
type GoogleSettings struct {
	ProjectID string

	Zone   string
	Region string

	VMType                     string
	SSHPublicKey               string
	EnableSSH                  bool
	VMPostDockerRunTargetState string // should become: SHUTDOWN | DELETE | KEEP

	DataStoreCollection string
	TaskPubSubTopic     string
}

const (
	// FlagProject ...
	FlagProject = "project"
	// FlagEntryPoint ...
	FlagEntryPoint = "entrypoint"
	// FlagVMType ...
	FlagVMType = "vm-type"
	// FlagZone defines zone to run VMs in
	FlagZone = "zone"
	// FlagRegion defines region of CFNs
	FlagRegion = "region"
	// FlagDetached sets dockers -d flag
	FlagDetached = "detached"
	// FlagEntrypoint IS NOT USED YET
	FlagEntrypoint = "entrypoint"

	// FlagSSHPublicKey can be deployed on VM instance
	FlagSSHPublicKey = "ssh-public-key"
	// FlagNoSSH disables inclusion of local user's SSH public keys in cloudInit
	FlagNoSSH = "no-ssh"

	// FireStoreCollection is the name of our firestore collection
	FireStoreCollection = "cloud-vm-docker-task"

	// TopicNameTaskQueue .. tbd: option/flag
	TopicNameTaskQueue = "cloud-vm-docker-task-queue"
)

// EnvironmentToGoogleSettings translates environment variables into a GoogleSettings struct.
func EnvironmentToGoogleSettings() GoogleSettings {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("CVD")
	s := GoogleSettings{
		ProjectID:                  viper.GetString(FlagProject),
		Zone:                       viper.GetString(FlagZone),
		Region:                     viper.GetString(FlagRegion),
		VMType:                     viper.GetString(FlagVMType),
		SSHPublicKey:               viper.GetString(FlagSSHPublicKey),
		EnableSSH:                  true,                // fixme
		VMPostDockerRunTargetState: "",                  // notyet
		DataStoreCollection:        FireStoreCollection, // fixme: static for now
		TaskPubSubTopic:            TopicNameTaskQueue,  // fixme: static for now
	}
	return s
}
