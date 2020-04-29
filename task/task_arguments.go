package task

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
)

// TaskArguments describes a Docker command to run on a specific type of VM
type TaskArguments struct {
	Image      string
	Command    []string
	EntryPoint string
	VMType     string
	VMID       string // we must pass this in upon creation request ... FIXME - Relict from P*bSub times?
	Subnet     string
	Tags       string
	// SSHPubKeys string holding one or more \n-separated ssh pubkeys
	SSHPubKeys string
}

// NewTaskArgumentsFromArgs returns a new TaskArguments based on CLI args
func NewTaskArgumentsFromArgs(image string, command []string, entryPoint, vmType, subnet, tags, sshKeys string) *TaskArguments {
	vmID := generateVMID()
	return &TaskArguments{
		Image:      image,
		Command:    command,
		EntryPoint: entryPoint,
		VMType:     vmType,
		VMID:       vmID,
		Subnet:     subnet,
		Tags:       tags,
		SSHPubKeys: sshKeys,
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

func generateVMID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("t%x", b[0:5])
}
