package cloud

import (
	"log"

	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// Setup creates necessary GoogleCloud infra for cloud-vm-docker operations --- NOT USABLE YET
func Setup(projectID string) {
	log.Print("SETTING UP infrastructure for cloud-vm-docker ...")
	bailOnError(createPubSubTopic(projectID, settings.TopicNameTaskQueue), "creating task queue")
	log.Print("SUCCESS setting up cloud-vm-docker cloud infra. Have fun!")
}

// Destroy removes infra created by setup routine
func Destroy(projectID string) {
	log.Print("DESTROYING infrastructure for cloud-vm-docker ...")
	bailOnError(deletePubSubTopic(projectID, settings.TopicNameTaskQueue), "deleting task queue")
	log.Print("SUCCESS destroying cloud-vm-docker cloud infra. Have fun!")
}

func bailOnError(err error, message string) {
	if err == nil {
		log.Printf("Success: %s", message)
	} else {
		log.Fatalf("ERROR: %s: %s", message, err.Error())
	}
}
