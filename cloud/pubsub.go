package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type CloudTask struct {
	Image      string
	Command    string
	EntryPoint string
	VMType     string
}

func NewCloudTaskFromArgs(image string, command string, entryPoint string, vmType string) *CloudTask {
	return &CloudTask{
		Image:      image,
		Command:    command,
		EntryPoint: entryPoint,
		VMType:     vmType,
	}
}

func NewCloudTaskFromBytes(data []byte) *CloudTask {
	task := CloudTask{}
	err := json.Unmarshal(data, task)
	if err != nil {
		log.Fatalf("Oooops, JSON task decoding error: %v", err)
	}
	return &task
}

func PubSubPushTask(task *CloudTask, projectID, topicID string) error {
	log.Printf("PUBLISHING TASK: %v to project %s topic %s", task, projectID, topicID)
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
	res := client.Topic(topicID).Publish(ctx, &m)
	log.Printf("PUBLISH res: %v", res)
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
