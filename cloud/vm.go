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

// fixme: this should receive a task ...
func CreateVM(projectID, zone, vmType, instanceName string) {
	log.Printf("Creating VM named %s of type %s in zone %s for project %s", instanceName, vmType, zone, projectID)
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// https://cloud.google.com/compute/docs/regions-zones#available
	machineTypeFQDN := fmt.Sprintf("zones/%s/machineTypes/%s",zone, vmType )
	prefix := "https://www.googleapis.com/compute/v1/projects/" + projectID
	rb := &compute.Instance{
		Name: instanceName,
		MachineType: machineTypeFQDN,
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskName:    "my-root-pd",
					SourceImage: getCOSImageLink(),
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
				Network: prefix + "/global/networks/default",
			},
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: "default",
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
				},
			},
		},
	}

	resp, err := computeService.Instances.Insert(projectID, zone, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// todo: if wait, then
	// https://github.com/googleapis/google-api-go-client/blob/master/examples/operation_progress.go

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("OK RESPONSE:\n%#v\n", resp)
}

func getCOSImageLink() string {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	image, err := computeService.Images.GetFromFamily("cos-cloud", "cos-stable").Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Selected VM Disk Image: %s (%s)", image.Name, image.Description)
	log.Printf("VM Disk Image SelfLink: %v", image.SelfLink)
	return image.SelfLink
}