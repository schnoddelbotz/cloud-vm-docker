package cloudfunctions

import (
	"context"
	"log"
	"net/http"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/cmd"
	"github.com/schnoddelbotz/cloud-vm-docker/handlers"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

var runtimeEnvironment *handlers.Environment

func init() {
	//var err error
	googleSettings := settings.EnvironmentToGoogleSettings()

	// gcloud / stackdriver logs have own timestamps, so drop Go's
	log.SetFlags(0)
	log.Printf("gotsfn %s package gotsfn-cfn init() with google settings from env: %v",
		cmd.AppVersion, googleSettings)

	runtimeEnvironment = handlers.NewEnvironment(googleSettings, true, true, true)
	log.Printf("Initialized runtime env: %v", runtimeEnvironment)
	// Google CloudFunction Go runtime requires us to use globals to share
	// resources between requests -- here: our bucket handle ...
	// https://cloud.google.com/functions/docs/concepts/go-runtime
	//handle, err = client.GetBucketHandle(settings)
	//if err != nil {
	//	log.Fatalf("gotsfn-cfn init() failed on GetBucketHandle(): %s", err)
	//}
	// FIXME: Maybe have relevant google clients (compute, pubsub, dataStore) as global?
}

// CloudVMDocker handles VMCreate, TaskStatus and TaskProgress requests
func CloudVMDocker(w http.ResponseWriter, r *http.Request) {
	handlers.CloudVMDocker(w, r) // , handle
}

// CloudVMDockerProcessor consumes a Pub/Sub message.
func CloudVMDockerProcessor(ctx context.Context, m cloud.PubSubMessage) error {
	return handlers.CloudVMDockerProcessor(ctx, m, runtimeEnvironment)
}
