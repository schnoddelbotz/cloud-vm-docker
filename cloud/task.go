package cloud

import "time"

// Task describes structure of our FireStore/DataStore documents
type Task struct {
	// TaskArguments hold `ctzz run` CLI arguments
	TaskArguments TaskArguments
	// Status tracks VM status
	Status string
	// VMID holds name of VM
	VMID string
	// ShutdownToken is known to VM itself, so it can request clean self-destruction
	ShutdownToken string
	// CreatedAt ...
	CreatedAt time.Time
}

const (
	TaskStatusCreated           = "created"           // initial state after creating new DataStore entry
	TaskStatusStarted           = "started"           // VM was created,
	TaskStatusRunning           = "running"           // client optionally reports this from within container + progress / eta
	TaskStatusKilled            = "killed"            // forceful `ctzz task-vm kill <VMID>`
	TaskStatusTerminatedFailure = "completed-failure" // container exited, code <> 0
	TaskStatusTerminatedSuccess = "completed-success" // container exited, code == 0
)
