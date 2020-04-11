package handlers

import (
	"context"
	"log"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
)

// CloudTaskZipZapProcessor consumes a Pub/Sub message.
func CloudTaskZipZapProcessor(ctx context.Context, m cloud.PubSubMessage) error {
	log.Printf("Request ctx: %v", ctx)
	task := cloud.NewCloudTaskFromBytes(m.Data)
	log.Printf("TASK: %v!", task)
	return nil
}
