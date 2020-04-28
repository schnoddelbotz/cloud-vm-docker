package handlers

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/compute/v1"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// Environment enables resource sharing between CFN requests, holds env settings + svc conns
type Environment struct {
	GoogleSettings  settings.GoogleSettings
	Context         context.Context
	FireStoreClient *firestore.Client
	ComputeService  *compute.Service
}

// NewEnvironment creates google service clients as requested
func NewEnvironment(googleSettings settings.GoogleSettings, withFireStoreClient bool, withComputeService bool) *Environment {
	var dataStoreClient *firestore.Client
	var computeService *compute.Service

	if googleSettings.AccessToken == "" {
		log.Fatal("NewEnvironment() FATAL: No CVD_TOKEN defined in environment. Cowardly refusing to expose service.")
	}
	if withFireStoreClient {
		dataStoreClient = cloud.NewFireStoreClient(context.Background(), googleSettings.ProjectID)
	}
	if withComputeService {
		computeService, _ = cloud.NewComputeService()
	}

	return &Environment{
		GoogleSettings:  googleSettings,
		Context:         context.Background(),
		FireStoreClient: dataStoreClient,
		ComputeService:  computeService,
	}
}
