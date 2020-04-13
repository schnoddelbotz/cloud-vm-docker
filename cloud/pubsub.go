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

// PubSubPushTask publishes a processing TaskArguments to given topic
func PubSubPushTask(task *TaskArguments, projectID, topicID string) error {
	log.Printf("PUBLISHING TASK to project %s topic %s:", projectID, topicID)
	log.Printf("  image  : %s", task.Image)
	log.Printf("  command: %q", task.Command)
	log.Printf("  vm_type: %s", task.VMType)
	log.Printf("  vm_id: %s", task.VMID)
	client, ctx := NewPubSubClient(projectID)
	data, _ := json.Marshal(task)
	m := pubsub.Message{
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

// NewPubSubClient returns client and its context, exits fatally on error
func NewPubSubClient(projectID string) (*pubsub.Client, context.Context) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	return client, ctx
}

func createPubSubTopic(projectID, topicID string) error {
	log.Printf("Creating pubsub topic '%s' in project %s ...", topicID, projectID)
	client, ctx := NewPubSubClient(projectID)
	t, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		return fmt.Errorf("CreateTopic: %v", err)
	}
	log.Printf("Topic created: %v\n", t)
	return nil
}

func deletePubSubTopic(projectID, topicID string) error {
	log.Printf("Deleting pubsub topic '%s' in project %s ... huh, dunno how, FXIME!", topicID, projectID)
	client, ctx := NewPubSubClient(projectID)
	topic := client.Topic(topicID)
	return topic.Delete(ctx)
}
