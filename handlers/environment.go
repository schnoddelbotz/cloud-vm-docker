package handlers

import (
	"context"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/compute/v1"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// Environment enables resource sharing between CFN requests, holds env settings + svc conns
type Environment struct {
	GoogleSettings  settings.GoogleSettings
	Context         context.Context
	PubSubClient    *pubsub.Client
	DataStoreClient *datastore.Client
	ComputeService  *compute.Service
}

// NewEnvironment creates google service clients as requested
func NewEnvironment(googleSettings settings.GoogleSettings, withPubSubClient bool, withDataStoreClient bool, withComputeService bool) *Environment {
	var pubSubClient *pubsub.Client
	var dataStoreClient *datastore.Client
	var computeService *compute.Service

	if withPubSubClient {
		pubSubClient, _ = cloud.NewPubSubClient(googleSettings.ProjectID)
	}
	if withDataStoreClient {
		dataStoreClient, _ = cloud.NewDataStoreClient(googleSettings.ProjectID)
	}
	if withComputeService {
		computeService, _ = cloud.NewComputeService()
	}

	return &Environment{
		GoogleSettings:  googleSettings,
		Context:         context.Background(),
		PubSubClient:    pubSubClient,
		DataStoreClient: dataStoreClient,
		ComputeService:  computeService,
	}
}
