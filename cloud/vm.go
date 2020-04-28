package cloud

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"google.golang.org/api/compute/v1"

	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// CreateVM spins up a ComputeEngine VM instance ...
func CreateVM(g settings.GoogleSettings, task Task) (*compute.Operation, error) {
	log.Printf("Creating VM named %s of type %s in zone %s in project %s", task.VMID, task.TaskArguments.VMType, g.Zone, g.ProjectID)
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
		if result.Status == "DONE" { // This "DONE" is NOT settings.TaskStatusDone!!
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

// DeleteInstanceByName ...
func DeleteInstanceByName(g settings.GoogleSettings, name string) error {
	log.Printf("DeleteInstanceByName called for: %s", name)
	computeClient, ctx := NewComputeService() // duh...
	_, err := computeClient.Instances.Delete(g.ProjectID, g.Zone, name).Context(ctx).Do()
	return err
}

func buildInstanceInsertionRequest(g settings.GoogleSettings, task Task) *compute.Instance {
	// instanceName, machineTypeFQDN, prefix, sshKeys, cloudInit string
	trueString := "true"
	machineTypeFQDN := fmt.Sprintf("zones/%s/machineTypes/%s", g.Zone, task.TaskArguments.VMType)
	prefix := "https://www.googleapis.com/compute/v1/projects/" + g.ProjectID
	cloudInit := buildCloudInit(g, task)
	var netIf compute.NetworkInterface
	if task.TaskArguments.Subnet == "" {
		netIf.AccessConfigs = []*compute.AccessConfig{
			{
				Type: "ONE_TO_ONE_NAT",
				Name: "External NAT",
			},
		}
		netIf.Network = prefix + "/global/networks/default"
	} else {
		subnetFormat := "projects/%s/regions/%s/subnetworks/%s"
		netIf.Subnetwork = fmt.Sprintf(subnetFormat, g.ProjectID, g.Region, task.TaskArguments.Subnet)
		log.Printf("Using non-default subnet for VM: %s", netIf.Subnetwork)
	}

	return &compute.Instance{
		Name:        task.VMID,
		MachineType: machineTypeFQDN,
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				// FIXME: make disk size and type a user choice
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskName:    "my-root-pd-" + task.VMID, // watch out! DiskName must be unique in project
					SourceImage: getCOSImageLink(),
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{&netIf},
		Tags:              nil, // TODO: --tags zrh-jump-imports,no-ip
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: "default", // FIXME make overridable
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring.write",
					"https://www.googleapis.com/auth/bigquery",
					"https://www.googleapis.com/auth/service.management.readonly",
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
	const tpl = `#cloud-config
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
    Environment="HOME=/home/cloudservice" "MGMT_TOKEN={{.Task.ManagementToken}}" "CVD_CFN_URL={{.ManagementURL}}" "CVD_VM_ID={{.Task.VMID}}"
    ExecStartPre=/usr/bin/docker-credential-gcr configure-docker
    ExecStartPre=/usr/bin/curl -s -XPOST -H"content-length: 0" -H"X-Authorization: ${MGMT_TOKEN}" ${CVD_CFN_URL}/status/${CVD_VM_ID}/booted
    ExecStart=/usr/bin/docker run \
        -v/var/run/docker.sock:/var/run/docker.sock \
        -v/home/cloudservice/.docker/config.json:/home/cloudservice/.docker/config.json \
        -eGCP_PROJECT={{.Google.ProjectID}} -eGCP_REGION={{.Google.Region}} \
        -eMGMT_TOKEN -eCVD_CFN_URL -eCVD_VM_ID \
        --name=cloud-vm-docker \
        {{.Task.TaskArguments.Image}} {{.QuotedCommand}}
    ExecStop=-/usr/bin/docker stop cloud-vm-docker
    ExecStopPost=/bin/sh -c "/usr/bin/docker inspect cloud-vm-docker --format='{''{'.Id'}''}' > /tmp/cid && /usr/bin/curl -s -d@/tmp/cid -H'X-Authorization: ${MGMT_TOKEN}' ${CVD_CFN_URL}/container/${CVD_VM_ID}/set"
    ExecStopPost=/usr/bin/sleep 15
    ExecStopPost=/usr/bin/curl -s -XPOST -H"content-length: 0" -H"X-Authorization: ${MGMT_TOKEN}" ${CVD_CFN_URL}/delete/${CVD_VM_ID}/${EXIT_STATUS}

runcmd:
- usermod -aG docker cloudservice
- docker-credential-gcr configure-docker
- systemctl daemon-reload
- systemctl start cloudservice.service
` // ATTN: Extra-careful to not put whitespace after \ for line concat above!

	// FIXME: The "ExecStopPost sleep" is only needed to "ensure" stackdriver logging agent
	//        container came up in time. Without it, logs would be lost on quickly failing containers.

	// ExecStopPost=/usr/bin/curl ... ${CVD_CFN_URL}/debug/${SERVICE_RESULT}/${EXIT_CODE}/${EXIT_STATUS}
	// --->
	// /debug/success/exited/0
	// /debug/exit-code/exited/127
	// /usr/bin/docker inspect cloud-vm-docker --format='{{.Id}}'

	// FIXME!!! There should be:
	// ExecStopPost=shutdown -h now ... as safety net
	// But then [Service] must run as root.

	type TemplateData struct {
		Google        settings.GoogleSettings
		Task          Task
		QuotedCommand string
		ManagementURL string
	}

	quotedCommand := ""
	if len(task.TaskArguments.Command) > 0 {
		// cloud-vm-docker task-vm create busybox sh -c 'sleep 3600'
		// --> busybox "sh" "-c" "sleep 3600"
		quotedCommand = strings.Trim(fmt.Sprintf("%q", task.TaskArguments.Command), "[]")
	}
	managementURL := fmt.Sprintf("https://%s-%s.cloudfunctions.net/CloudVMDocker", g.Region, g.ProjectID)
	templateData := TemplateData{Google: g, Task: task, QuotedCommand: quotedCommand, ManagementURL: managementURL}
	t := template.Must(template.New("cloud-init").Parse(tpl))

	var result bytes.Buffer
	err := t.Execute(&result, templateData)

	if err != nil {
		log.Fatalf("executing cloud-init template: %s", err)
	}

	return result.String()
}
