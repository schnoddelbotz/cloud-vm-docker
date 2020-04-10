package cloud

import (
	"log"

	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

func Setup(projectID string) {
	log.Print("SETTING UP infrastructure for cloud-task-zip-zap ...")
	bailOnError(createPubSubTopic(projectID, settings.TopicNameTaskQueue), "creating task queue")
	bailOnError(createPubSubTopic(projectID, settings.TopicNameStatusQueue), "creating status queue")
	log.Print("SUCCESS setting up cloud-task-zip-zap cloud infra. Have fun!")
}

func Destroy(projectID string) {
	log.Print("DESTROYING infrastructure for cloud-task-zip-zap ...")
	bailOnError(deletePubSubTopic(projectID, settings.TopicNameTaskQueue), "deleting task queue")
	bailOnError(deletePubSubTopic(projectID, settings.TopicNameStatusQueue), "deleting status queue")
	log.Print("SUCCESS destroying cloud-task-zip-zap cloud infra. Have fun!")
}

func bailOnError(err error, message string) {
	if err == nil {
		log.Printf("Success: %s", message)
	} else {
		log.Fatalf("ERROR: %s: %s", message, err.Error())
	}
}