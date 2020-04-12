package handlers

import (
	"context"
	"log"
	"os"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
)

// CloudTaskZipZapProcessor consumes a Pub/Sub message.
func CloudTaskZipZapProcessor(ctx context.Context, m cloud.PubSubMessage) error {
	log.Printf("Request ctx: %v", ctx)
	task := cloud.NewCloudTaskFromBytes(m.Data)
	project := os.Getenv("GCP_PROJECT")
	log.Printf("TASK: image='%s' command=%q vmtype='%s'!", task.Image, task.Command, task.VMType)
	cloud.StoreTask(project, *task)
	return nil
}
