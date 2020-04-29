package task

import (
	"time"
)

// Task describes structure of our FireStore/FireStore documents
type Task struct {
	// TaskArguments hold `cloud-vm-docker run` CLI arguments
	TaskArguments TaskArguments
	// Status tracks VM status
	Status string
	// VMID holds name of VM
	VMID string
	// InstanceID holds ID of GCE instance, only known after creation. It's a uint64, but datastore no supporty.
	InstanceID string
	// ManagementToken is known to VM itself, so it can request clean self-destruction
	ManagementToken string
	// CreatedAt ...
	CreatedAt time.Time

	// DockerExitCode from docker run command on VM
	DockerExitCode int
	// DockerContainerId stores container ID on VM, to enable StackDriver log filtering
	DockerContainerId string
}

const (
	// TaskStatusCreated : initial state after creating new FireStore entry
	TaskStatusCreated = "created"
	// TaskStatusBooted : VM booted, tries to run docker pull on our image now
	TaskStatusBooted = "booted"
	// TaskStatusRunning : client has pulled image and docker command should be running now
	TaskStatusRunning = "running"
	// TaskStatusDone : Container has exited, VM destruction requested
	TaskStatusDone = "done"
)
