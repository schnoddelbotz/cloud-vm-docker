package settings

// Settings configure a CLI or CloudFunction ctzz instance
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

	// ActionSubmit ...
	ActionSubmit = "submit"
	// ActionWait ...
	ActionWait = "wait"
	// ActionKill ...
	ActionKill = "kill"

	// FireStoreCollection is the name of our firestore collection
	FireStoreCollection = "ctzz-task"

	// TopicNameTaskQueue .. tbd: option/flag
	TopicNameTaskQueue = "ctzz-task-queue"
	// TopicNameStatusQueue ...
	TopicNameStatusQueue = "ctzz-status-queue"
)
