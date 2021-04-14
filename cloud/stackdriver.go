package cloud

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

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
	filter := fmt.Sprintf(`resource.type="gce_instance"
jsonPayload."cos.googleapis.com/container_id"="%s"`, containerID)
	baseLink := `https://console.cloud.google.com/logs/viewer?resource=gce_instance/instance_id/%d&project=%s&advancedFilter=%s`
	return fmt.Sprintf(baseLink, GCEInstanceID, project, url.QueryEscape(filter))
}

func PrintLogEntries(projectID, instanceID, containerID string, createdAt time.Time) {
	//log.Printf("Printing logs for %s/instanceID=%s/container=%s", projectID, instanceID, containerID)
	ctx := context.Background()
	client := NewStackDriverClient(ctx, projectID)

	// FIXME:
	logFilterFormat := `
		logName = "projects/%s/logs/cos_containers"
		resource.type="gce_instance"
		timestamp >= "%s"
		resource.labels.instance_id="%s"
		jsonPayload."cos.googleapis.com/container_id"="%s"`
	filter := fmt.Sprintf(logFilterFormat, projectID, createdAt.Format(time.RFC3339), instanceID, containerID)
	//log.Printf("Using filter:\n%s", filter)

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
