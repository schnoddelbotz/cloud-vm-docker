package cloud

import (
	"context"
	"fmt"
	"google.golang.org/api/compute/v1"
	"log"
	"time"
)

// https://github.com/googleapis/google-api-go-client/blob/master/examples/compute.go
// https://cloud.google.com/compute/docs/reference/rest/v1/instances/insert
// https://godoc.org/google.golang.org/api/compute/v1

// CreateVM spins up a ComputeEngine VM instance ...
// fixme: this should receive a task ...
//func CreateVM(projectID, zone, vmType, instanceName string) {
func CreateVM(projectID, zone string, task Task, sshKeys string) (*compute.Operation, error) {
	log.Printf("Creating VM named %s of type %s in zone %s for project %s", task.VMID, task.TaskArguments.VMType, zone, projectID)
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// https://cloud.google.com/compute/docs/regions-zones#available
	machineTypeFQDN := fmt.Sprintf("zones/%s/machineTypes/%s", zone, task.TaskArguments.VMType)
	prefix := "https://www.googleapis.com/compute/v1/projects/" + projectID
	cloudInit := buildCloudInit(projectID, "cfnRegion", task.TaskArguments.Image, task.TaskArguments.Command)
	rb := buildInstanceInsertionRequest(task.VMID, machineTypeFQDN, prefix, sshKeys, cloudInit)
	resp, err := computeService.Instances.Insert(projectID, zone, rb).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	// todo: if wait, then
	// https://github.com/googleapis/google-api-go-client/blob/master/examples/operation_progress.go

	// TODO: Change code below to process the `resp` object:
	// FIXME : OR JUST RETURN IT! or just SelfLink to operation
	// fmt.Printf("OK RESPONSE:\n%#v\n", resp)
	return resp, nil
}

func WaitForOperation(project, zone, operation string) {
	log.Printf("Waiting for operation %s in zone %s project %s", operation, zone, project)
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to create compute client. Should have re-used anyway (FIXME)")
	}
	// todo: add max / timeout
	time.Sleep(1*time.Second)
	waited := 1
	for {
		result, err := computeService.ZoneOperations.Get(project, zone, operation).Do()
		if err != nil {
			log.Fatalf("Error getting operations: %s", err)
		}
		if result.Status == "DONE" {
			log.Printf("FINALLY ... found DONE status after %d seconds", waited)
			return
		}
		if waited % 10 == 0 {
			log.Printf("Already waited %d seconds ...", waited)
		}
		time.Sleep(1*time.Second)
		waited += 1
	}

}

func buildInstanceInsertionRequest(instanceName, machineTypeFQDN, prefix, sshKeys, cloud_init string) *compute.Instance {
	true_string := "true"
	return &compute.Instance{
		Name:        instanceName,
		MachineType: machineTypeFQDN,
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskName:    "my-root-pd-"+instanceName,
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
		Metadata: &compute.Metadata{
			Items:           []*compute.MetadataItems{
				{
					Key: "ssh-keys",
					Value: &sshKeys,
				},
				{
					Key: "google-logging-enabled",
					Value: &true_string,
				},
				{
					Key: "user-data",
					Value: &cloud_init,
				},
			},
		},
	}
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
