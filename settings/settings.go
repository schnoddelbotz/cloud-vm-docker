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
	// FlagImage ...
	FlagImage = "image"
	// FlagCommand ...
	FlagCommand = "command"
	// FlagArgs ...
	FlagArgs = "args"
	// FlagArgsFile ...
	FlagArgsFile = "args-file"
	// FlagEntryPoint ...
	FlagEntryPoint = "entrypoint"
	// FlagVMType ...
	FlagVMType = "vm-type"

	// ActionSubmit ...
	ActionSubmit = "submit"
	// ActionWait ...
	ActionWait = "wait"
	// ActionKill ...
	ActionKill = "kill"

	// TopicNameTaskQueue .. tbd: option/flag
	TopicNameTaskQueue = "ctzz-task-queue"
	// TopicNameStatusQueue ...
	TopicNameStatusQueue = "ctzz-status-queue"
)
