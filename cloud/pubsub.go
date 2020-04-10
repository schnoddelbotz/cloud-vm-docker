package cloud

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func createPubSubTopic(projectID, topicID string) error {
	fmt.Printf("Creating pubsub topic '%s' in project %s ...\n", topicID, projectID)
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		return fmt.Errorf("CreateTopic: %v", err)
	}
	fmt.Printf("Topic created: %v\n", t)
	return nil
}
