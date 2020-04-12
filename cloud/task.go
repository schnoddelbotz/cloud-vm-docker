package cloud

import "time"

// Task describes structure of our FireStore/DataStore documents
type Task struct {
	TaskArguments TaskArguments
	Status        string
	VMID          string
	ShutdownToken string
	CreatedAt     time.Time
}
