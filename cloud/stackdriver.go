package cloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/iterator"
)

// GetLogLinkForVM returns a google web console link to VM's stackdriver logs
func GetLogLinkForVM(project string, GCEInstanceID uint64) string {
	baseLink := `https://console.cloud.google.com/logs/viewer?resource=gce_instance/instance_id/%d&project=%s`
	return fmt.Sprintf(baseLink, GCEInstanceID, project)
}

// GetLogLinkForContainer as GetLogLinkForVM, but filtered to specific Docker container
func GetLogLinkForContainer(project string, GCEInstanceID uint64, containerID string) string {
	baseLink := `https://console.cloud.google.com/logs/viewer?resource=gce_instance/instance_id/%d&project=%s&advancedFilter=resource.type="gce_instance"+jsonPayload.container_id="%s"`
	return fmt.Sprintf(baseLink, GCEInstanceID, project, containerID)
}

func PrintLogEntries(projectID, instanceID, containerID string) {
	//log.Printf("Printing logs for %s/instanceID=%s/container=%s", projectID, instanceID, containerID)
	ctx := context.Background()
	client := NewStackDriverClient(ctx, projectID)

	// FIXME:
	logFilterFormat := `
		logName = "projects/%s/logs/cos_containers" 
		jsonPayload.container_id="%s" 
		resource.type="gce_instance" 
		resource.labels.instance_id="%s"`
	filter := fmt.Sprintf(logFilterFormat, projectID, containerID, instanceID)

	it := client.Entries(ctx, logadmin.Filter(filter))
	var entries []*logging.Entry
	_, err := iterator.NewPager(it, 512, "").NextPage(&entries)
	if err != nil {
		log.Fatalf("problem getting logs: %v", err)
		return
	}
	for _, entry := range entries {
		message := entry.Payload.(*structpb.Struct).GetFields()["message"].GetStringValue() // RLY???
		fmt.Printf("%s %s\n", entry.Timestamp.Local().Format("2006-01-02T15:04:05"), strings.TrimRight(message, "\n"))
	}

	// TODO:
	// log.Printf("NOTE! WIP! Only showed last 512 log entries. There might be more. FIXME.")
}

func NewStackDriverClient(ctx context.Context, projectID string) *logadmin.Client {
	client, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error creating StackDriver client: %s", err)
	}
	return client
}
