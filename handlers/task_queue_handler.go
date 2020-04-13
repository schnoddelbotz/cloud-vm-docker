package handlers

import (
	"context"
	"log"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
)

// CloudVMDockerProcessor consumes a Pub/Sub message.
func CloudVMDockerProcessor(_ context.Context, m cloud.PubSubMessage, runtimeEnvironment *Environment) error {
	task := cloud.NewCloudTaskArgsFromBytes(m.Data)
	g := runtimeEnvironment.GoogleSettings
	log.Printf("TASK: project='%s' image='%s' command=%q vmtype='%s'!", g.ProjectID, task.Image, task.Command, task.VMType)
	cloud.StoreNewTask(g.ProjectID, *task)
	log.Printf("Created task successfully, should now spawn VM... FIXME ... and return VM id :-/")
	return nil
}
