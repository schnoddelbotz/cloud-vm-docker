package cloud

import (
	"encoding/json"
	"log"
)

// TaskArguments describes a Docker command to run on a specific type of VM
type TaskArguments struct {
	Image      string
	Command    []string
	EntryPoint string
	VMType     string
}

// NewCloudTaskArgsFromArgs returns a new TaskArguments based on CLI args
func NewCloudTaskArgsFromArgs(image string, command []string, entryPoint string, vmType string) *TaskArguments {
	return &TaskArguments{
		Image:      image,
		Command:    command,
		EntryPoint: entryPoint,
		VMType:     vmType,
	}
}

// NewCloudTaskArgsFromBytes returns a new TaskArguments based on a (pubsub) JSON message
func NewCloudTaskArgsFromBytes(data []byte) *TaskArguments {
	task := TaskArguments{}
	err := json.Unmarshal(data, &task)
	if err != nil {
		log.Fatalf("Oooops, JSON task decoding error: %v", err)
	}
	return &task
}
