package cloud

import (
	"context"
	"fmt"
	"github.com/schnoddelbotz/cloud-vm-docker/settings"
	"log"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
)

// CreateVM spins up a ComputeEngine VM instance ...
func CreateVM(g settings.GoogleSettings, task Task, sshKeys string) (*compute.Operation, error) {
	log.Printf("Creating VM named %s of type %s in zone %s for project %s", task.VMID, task.TaskArguments.VMType, g.Zone, g.ProjectID)
	computeService, ctx := NewComputeService()

	rb := buildInstanceInsertionRequest(g, task)
	resp, err := computeService.Instances.Insert(g.ProjectID, g.Zone, rb).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WaitForOperation guess what
func WaitForOperation(project, zone, operation string) {
	log.Printf("Waiting for operation %s in zone %s project %s", operation, zone, project)
	computeService, _ := NewComputeService()
	// todo: add max / timeout
	time.Sleep(1 * time.Second)
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
		if waited%10 == 0 {
			log.Printf("Already waited %d seconds ...", waited)
		}
		time.Sleep(1 * time.Second)
		waited++
	}
}

// NewComputeService returns a compute service client and its context; fatally fails on error
func NewComputeService() (*compute.Service, context.Context) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to create compute client. Should have re-used anyway (FIXME)")
	}
	return computeService, ctx
}

func buildInstanceInsertionRequest(g settings.GoogleSettings, task Task) *compute.Instance {
	// instanceName, machineTypeFQDN, prefix, sshKeys, cloudInit string
	trueString := "true"
	machineTypeFQDN := fmt.Sprintf("zones/%s/machineTypes/%s", g.Zone, task.TaskArguments.VMType)
	prefix := "https://www.googleapis.com/compute/v1/projects/" + g.ProjectID
	cloudInit := buildCloudInit(g, task)
	return &compute.Instance{
		Name:        task.VMID,
		MachineType: machineTypeFQDN,
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskName:    "my-root-pd-" + task.VMID, // watch out! DiskName must be unique in project
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
			Items: []*compute.MetadataItems{
				{
					Key:   "ssh-keys",
					Value: &task.SSHPubKeys,
				},
				{
					Key:   "google-logging-enabled",
					Value: &trueString,
				},
				{
					Key:   "user-data",
					Value: &cloudInit,
				},
			},
		},
	}
}

func getCOSImageLink() string {
	computeService, _ := NewComputeService()
	image, err := computeService.Images.GetFromFamily("cos-cloud", "cos-stable").Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Selected VM Disk Image: %s (%s)", image.Name, image.Description)
	log.Printf("VM Disk Image SelfLink: %v", image.SelfLink)
	return image.SelfLink
}

func buildCloudInit(g settings.GoogleSettings, task Task) string {
	project := g.ProjectID
	cfnRegion := g.Region
	image := task.TaskArguments.Image
	command := task.TaskArguments.Command
	// FIXME use text/template
	//func buildCloudInit(project, cfnRegion, image string, command []string) string {
	// should set vm shutdown token
	// should use task as first arg?
	// should quote all command parts
	myCommand := strings.Join(command, " ")
	return fmt.Sprintf(`#cloud-config
users:
- name: cloudservice
  uid: 2000

write_files:
- path: /etc/systemd/system/cloudservice.service
  permissions: 0644
  owner: root
  content: |
    [Unit]
    Description=Start inventory-optimisation docker container
    Wants=gcr-online.target
    After=gcr-online.target

    [Service]
    User=cloudservice
    Restart=no
    Environment="HOME=/home/cloudservice"
    ExecStartPre=/usr/bin/docker-credential-gcr configure-docker
    ExecStart=/usr/bin/docker run --rm \
        -v/var/run/docker.sock:/var/run/docker.sock \
        -v/home/cloudservice/.docker/config.json:/home/cloudservice/.docker/config.json \
        -eGCP_PROJECT=%s -eGCP_REGION=%s --name=cloud-vm-docker %s %s
    ExecStop=/usr/bin/docker stop cloud-vm-docker
    ExecStopPost=/usr/bin/docker rm cloud-vm-docker

runcmd:
- usermod -aG docker cloudservice
- docker-credential-gcr configure-docker
- systemctl daemon-reload
- systemctl start cloudservice.service
`, project, cfnRegion, image, myCommand)
}
