package settings

// Settings configure a CLI or CloudFunction ctzz instance
type Settings struct {
	PrintVersion  bool
	PrintRawJSON  bool
	CopyClipboard bool
	NoOutput      bool
	RawOutput     bool
	Provider      string
	File          string
	Server        string
	Google        GoogleCloudSettings
}

// GoogleCloudSettings define anything Google related (project, service account, ...)
type GoogleCloudSettings struct {
	ProjectID           string
	DataStoreCollection string
}

const (
	// ActionSubmit
	ActionSubmit = "submit"
	// ActionWait
	ActionWait = "wait"
	// ActionKill
	ActionKill = "kill"
)
