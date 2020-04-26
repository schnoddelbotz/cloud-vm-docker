package cloud

import "fmt"

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
