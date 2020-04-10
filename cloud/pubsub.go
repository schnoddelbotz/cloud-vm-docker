package cloud

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

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
