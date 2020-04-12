package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/pubsub/pubsub.go

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Task describes a Docker command to run on a specific type of VM
type Task struct {
	Image      string
	Command    []string
	EntryPoint string
	VMType     string
}

// NewCloudTaskFromArgs returns a new Task based on CLI args
func NewCloudTaskFromArgs(image string, command []string, entryPoint string, vmType string) *Task {
	return &Task{
		Image:      image,
		Command:    command,
		EntryPoint: entryPoint,
		VMType:     vmType,
	}
}

// NewCloudTaskFromBytes returns a new Task based on a (pubsub) JSON message
func NewCloudTaskFromBytes(data []byte) *Task {
	task := Task{}
	err := json.Unmarshal(data, &task)
	if err != nil {
		log.Fatalf("Oooops, JSON task decoding error: %v", err)
	}
	return &task
}

// PubSubPushTask publishes a processing Task to given topic
func PubSubPushTask(task *Task, projectID, topicID string) error {
	log.Printf("PUBLISHING TASK to project %s topic %s:", projectID, topicID)
	log.Printf("  image  : %s", task.Image)
	log.Printf("  command: %q", task.Command)
	log.Printf("  vm_type: %s", task.VMType)
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	data, _ := json.Marshal(task)
	m := pubsub.Message{
		ID:          "1234",
		Data:        data,
		Attributes:  nil,
		PublishTime: time.Time{},
	}
	if _, err := client.Topic(topicID).Publish(ctx, &m).Get(ctx); err != nil {
		log.Fatalf("FAILED to publish: %v", err)
	}
	log.Println("PUBLISH successs")
	return nil
}

func createPubSubTopic(projectID, topicID string) error {
	log.Printf("Creating pubsub topic '%s' in project %s ...", topicID, projectID)
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		return fmt.Errorf("CreateTopic: %v", err)
	}
	log.Printf("Topic created: %v\n", t)
	return nil
}

func deletePubSubTopic(projectID, topicID string) error {
	log.Printf("Deleting pubsub topic '%s' in project %s ... huh, dunno how, FXIME!", topicID, projectID)
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	topic := client.Topic(topicID)
	return topic.Delete(ctx)
}
