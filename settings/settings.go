package settings

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// GoogleSettings define anything Google related (project, service account, ...)
type GoogleSettings struct {
	ProjectID                  string
	Zone                       string
	Region                     string
	VMType                     string
	SSHPublicKey               string
	DisableSSH                 bool
	VMPostDockerRunTargetState string // should become: SHUTDOWN | DELETE | KEEP
	DataStoreCollection        string
	AccessToken                string // access protects the HTTP CFN
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
	// FlagWait tells run command to wait until container completed
	FlagWait = "wait"
	// FlagToken is the access token for CloudVMDocker HTTP CloudFunction for run/ps/...
	FlagToken = "token"

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
func EnvironmentToGoogleSettings(permitEmptyToken bool) GoogleSettings {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("CVD")
	accessToken := viper.GetString(FlagToken)
	if accessToken == "" {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		accessToken = fmt.Sprintf("t%x", b[0:5])
		if !permitEmptyToken {
			log.Fatalf("ERROR! Empty CVD_TOKEN not permitted. Define via env or --token ...")
		}
		log.Printf("Warning! Empty CVD_TOKEN -- generated random one for use: %s", accessToken)
	}
	s := GoogleSettings{
		ProjectID:                  viper.GetString(FlagProject),
		Zone:                       viper.GetString(FlagZone),
		Region:                     viper.GetString(FlagRegion),
		VMType:                     viper.GetString(FlagVMType),
		AccessToken:                accessToken,
		SSHPublicKey:               viper.GetString(FlagSSHPublicKey),
		DisableSSH:                 viper.GetBool(FlagNoSSH),
		VMPostDockerRunTargetState: "",                  // notyet
		DataStoreCollection:        FireStoreCollection, // fixme: static for now
	}
	return s
}
