package settings

// Settings configure a CLI or CloudFunction cloud-vm-docker instance
type Settings struct {
	NoOutput  bool
	RawOutput bool
	Provider  string
	File      string
	Server    string
	Google    GoogleCloudSettings
}

// GoogleCloudSettings define anything Google related (project, service account, ...)
type GoogleCloudSettings struct {
	ProjectID           string
	DataStoreCollection string
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

	// FlagSSHPublicKey can be deployed on VM instance
	FlagSSHPublicKey = "ssh-public-key"
	// FlagNoSSH disables inclusion of local user's SSH public keys in cloudInit
	FlagNoSSH = "no-ssh"

	// ActionSubmit ...
	ActionSubmit = "submit"
	// ActionWait ...
	ActionWait = "wait"
	// ActionKill ...
	ActionKill = "kill"

	// FireStoreCollection is the name of our firestore collection
	FireStoreCollection = "cloud-vm-docker-task"

	// TopicNameTaskQueue .. tbd: option/flag
	TopicNameTaskQueue = "cloud-vm-docker-task-queue"
	// TopicNameStatusQueue ...
	TopicNameStatusQueue = "cloud-vm-docker-status-queue"
)
