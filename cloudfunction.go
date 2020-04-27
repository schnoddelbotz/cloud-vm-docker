package cloudfunction

import (
	"log"
	"net/http"

	"github.com/schnoddelbotz/cloud-vm-docker/handlers"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
	// to load viper defaults for our flags...
	_ "github.com/schnoddelbotz/cloud-vm-docker/cmd"
)

var runtimeEnvironment *handlers.Environment

func init() {
	// gcloud / stackdriver logs have own timestamps, so drop Go's
	log.SetFlags(0)

	// import environment vars, using same defaults as CLI
	googleSettings := settings.EnvironmentToGoogleSettings(true)
	log.Printf("CLOUD-VM-DOCKER initialized with settings from env: %v", googleSettings)

	// we initialize all clients here, albeit different needs of CFNs. Solve.
	runtimeEnvironment = handlers.NewEnvironment(googleSettings, true, true)
	log.Printf("Initialized runtime env: %v", runtimeEnvironment)
}

// CloudVMDocker handles VMCreate, TaskStatus and TaskProgress requests
func CloudVMDocker(w http.ResponseWriter, r *http.Request) {
	handlers.CloudVMDocker(w, r, runtimeEnvironment)
}
