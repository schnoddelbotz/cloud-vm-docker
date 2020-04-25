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
	VMID       string // we must pass this in upon creation request ... FIXME - Relict from P*bSub times?
}

// NewTaskArgumentsFromArgs returns a new TaskArguments based on CLI args
func NewTaskArgumentsFromArgs(image string, command []string, entryPoint string, vmType string) *TaskArguments {
	vmID := generateVMID()
	return &TaskArguments{
		Image:      image,
		Command:    command,
		EntryPoint: entryPoint,
		VMType:     vmType,
		VMID:       vmID,
	}
}

// NewTaskArgumentsFromBytes returns a new TaskArguments based on a JSON message
func NewTaskArgumentsFromBytes(data []byte) *TaskArguments {
	task := TaskArguments{}
	err := json.Unmarshal(data, &task)
	if err != nil {
		log.Fatalf("Oooops, JSON task decoding error: %v", err)
	}
	return &task
}
