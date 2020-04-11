package cloudfunctions

import (
	"context"
	"github.com/schnoddelbotz/cloud-task-zip-zap/cloud"
	"log"
	"net/http"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cmd"
	"github.com/schnoddelbotz/cloud-task-zip-zap/handlers"
)

func init() {
	//var err error
	settings := cmd.EnvironmentToSettings()

	// gcloud / stackdriver logs have own timestamps, so drop Go's
	log.SetFlags(0)
	log.Printf("gotsfn %s package gotsfn-cfn init() with settings: %v",
		cmd.AppVersion, settings)

	// Google CloudFunction Go runtime requires us to use globals to share
	// resources between requests -- here: our bucket handle ...
	// https://cloud.google.com/functions/docs/concepts/go-runtime
	//handle, err = client.GetBucketHandle(settings)
	//if err != nil {
	//	log.Fatalf("gotsfn-cfn init() failed on GetBucketHandle(): %s", err)
	//}
}

// CloudTaskZipZap handler manages secrets but also delivers web app.
// It is deployed as single Google CloudFunction.
func CloudTaskZipZap(w http.ResponseWriter, r *http.Request) {
	handlers.CloudTaskZipZap(w, r) // , handle
}

// CloudTaskZipZapProcessor consumes a Pub/Sub message.
func CloudTaskZipZapProcessor(ctx context.Context, m cloud.PubSubMessage) error {
	return handlers.CloudTaskZipZapProcessor(ctx, m)
}
