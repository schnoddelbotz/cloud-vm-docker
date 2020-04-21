package cloud

import "time"

// Task describes structure of our FireStore/DataStore documents
type Task struct {
	// TaskArguments hold `cloud-vm-docker run` CLI arguments
	TaskArguments TaskArguments
	// Status tracks VM status
	Status string
	// VMID holds name of VM
	VMID string
	// ManagementToken is known to VM itself, so it can request clean self-destruction
	ManagementToken string
	// CreatedAt ...
	CreatedAt time.Time
	// SSHPubKeys string holding one or more \n-separated ssh pubkeys
	SSHPubKeys string
	// DockerExitCode from docker run command on VM
	DockerExitCode int
}

const (
	// TaskStatusCreated : initial state after creating new DataStore entry
	TaskStatusCreated = "created"
	// TaskStatusStarted : VM was created
	TaskStatusStarted = "started"
	// TaskStatusRunning : client optionally reports this from within container + progress / eta
	TaskStatusRunning = "running"
	// TaskStatusKilled : forceful `cloud-vm-docker task-vm kill <VMID>`
	TaskStatusKilled = "killed"
	// TaskStatusTerminatedFailure : container exited, code <> 0
	TaskStatusTerminatedFailure = "completed-failure"
	// TaskStatusTerminatedSuccess : container exited, code == 0
	TaskStatusTerminatedSuccess = "completed-success"
)
