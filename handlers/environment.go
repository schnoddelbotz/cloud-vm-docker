package handlers

import (
	"context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/compute/v1"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// Environment enables resource sharing between CFN requests, holds env settings + svc conns
type Environment struct {
	GoogleSettings  settings.GoogleSettings
	Context         context.Context
	DataStoreClient *datastore.Client
	ComputeService  *compute.Service
}

// NewEnvironment creates google service clients as requested
func NewEnvironment(googleSettings settings.GoogleSettings, withDataStoreClient bool, withComputeService bool) *Environment {
	var dataStoreClient *datastore.Client
	var computeService *compute.Service

	eContext := context.Background()

	if withDataStoreClient {
		dataStoreClient = cloud.NewDataStoreClient(eContext, googleSettings.ProjectID)
	}
	if withComputeService {
		computeService, _ = cloud.NewComputeService()
	}

	return &Environment{
		GoogleSettings:  googleSettings,
		Context:         eContext,
		DataStoreClient: dataStoreClient,
		ComputeService:  computeService,
	}
}
