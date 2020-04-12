package cloud

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/compute/v1"
)

// https://github.com/googleapis/google-api-go-client/blob/master/examples/compute.go
// https://cloud.google.com/compute/docs/reference/rest/v1/instances/insert
// https://godoc.org/google.golang.org/api/compute/v1

func CreateVM(projectID, vmType, instanceName string) {
	log.Printf("Creating VM named %s of type %s in project %s", instanceName, vmType, projectID)
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// The name of the zone for this request.
	zone := "europe-west1-b" // TODO: Update placeholder value.
	// https://cloud.google.com/compute/docs/regions-zones#available

	rb := &compute.Instance{
		// TODO: Add desired fields of the request body.
	}

	resp, err := computeService.Instances.Insert(projectID, zone, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("OK RESPONSE:\n%#v\n", resp)
}