package cloud

import "log"

// https://github.com/googleapis/google-api-go-client/blob/master/examples/compute.go

func CreateVM(projectID, vmType string, vmName string) {
	log.Printf("Creating VM named %s of type %s in project %s", vmName, vmType, projectID)
}