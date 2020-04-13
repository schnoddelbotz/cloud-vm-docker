package handlers

import (
	"context"
	"log"
	"os"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
)

// CloudVMDockerProcessor consumes a Pub/Sub message.
func CloudVMDockerProcessor(_ context.Context, m cloud.PubSubMessage, runtimeEnvironment *Environment) error {
	task := cloud.NewCloudTaskArgsFromBytes(m.Data)
	project := os.Getenv("CVD_PROJECT")
	log.Printf("TASK: project='%s' image='%s' command=%q vmtype='%s'!", project, task.Image, task.Command, task.VMType)
	log.Printf("MY RTENV: %v", runtimeEnvironment)
	cloud.StoreNewTask(project, *task)
	log.Printf("Created task successfully, should now spawn VM... FIXME")
	return nil
}
