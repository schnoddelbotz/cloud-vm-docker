package settings

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"

	"github.com/schnoddelbotz/cloud-vm-docker/task"
)

// RuntimeSettings define anything Google related (project, service account, ...)
type RuntimeSettings struct {
	ProjectID string
	//Zone      string
	Region         string
	TaskArgs       task.TaskArguments
	NoSSH          bool   // not needed to be saved -- if disabled, there are no keys.
	Token          string // access protects the HTTP CFN
	ServiceAccount string
	Detached       bool
	Wait           bool
	Verbose        bool
	PrintLogs      bool
}

const (
	// FlagProject ...
	FlagProject = "project"
	// FlagZone defines zone to run VMs in
	FlagZone = "zone"
	// FlagRegion defines region of CFNs
	FlagRegion = "region"

	// FlagDetached sets dockers -d flag (for runCmd)
	FlagDetached = "detached"
	// FlagEntrypoint IS NOT USED YET
	FlagEntrypoint = "entrypoint"
	// FlagWait tells run command to wait until container completed
	FlagWait = "wait"
	// FlagToken is the access token for CloudVMDocker HTTP CloudFunction for run/ps/...
	FlagToken = "token"
	// FlagVerbose is just a bool here
	FlagVerbose = "verbose"
	// FlagPrintLogs en/disables printing logs after waiting for VM/Docker completion
	FlagPrintLogs = "print-logs"

	// TaskArgs
	// FlagVMType selects VM type from https://cloud.google.com/compute/docs/machine-types
	FlagVMType = "vm-type"
	// FlagSubnet is used for VM creation
	FlagSubnet = "subnet"
	// FlagTags defines tags to apply to VM creation (comma-separated)
	FlagTags = "tags"
	// FlagEntryPoint NOT YET
	FlagEntryPoint     = "entrypoint"
	FlagServiceAccount = "service-account"

	// FlagSSHPublicKey can be deployed on VM instance
	FlagSSHPublicKey = "ssh-public-key"
	// FlagNoSSH disables inclusion of local user's SSH public keys in cloudInit
	FlagNoSSH = "no-ssh"

	// FireStoreCollection is the name of our firestore collection (static for now)
	FireStoreCollection = "cloud-vm-docker-task"
)

// ViperToRuntimeSettings translates environment variables into a RuntimeSettings struct.
func ViperToRuntimeSettings(permitEmptyToken bool) RuntimeSettings {
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
	s := RuntimeSettings{
		ProjectID: viper.GetString(FlagProject),
		//Zone:      viper.GetString(FlagZone),
		Region:         viper.GetString(FlagRegion),
		Token:          accessToken,
		ServiceAccount: viper.GetString(FlagServiceAccount),
		NoSSH:          viper.GetBool(FlagNoSSH),
		Detached:       viper.GetBool(FlagDetached),
		Wait:           viper.GetBool(FlagWait),
		Verbose:        viper.GetBool(FlagVerbose),
		PrintLogs:      viper.GetBool(FlagPrintLogs),
		TaskArgs: *task.NewTaskArgumentsFromArgs("", nil, viper.GetString(FlagEntryPoint),
			viper.GetString(FlagVMType), viper.GetString(FlagZone), viper.GetString(FlagSubnet),
			viper.GetString(FlagTags), viper.GetString(FlagSSHPublicKey)),
	}
	return s
}
