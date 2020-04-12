package cloud

// https://cloud.google.com/datastore/docs/reference/libraries

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
)

type Document struct {
	Task Task
	Status string
}

func StoreTask(projectID string, task Task) {
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the kind for the new entity. // CollectionName
	kind := "Task"
	// Sets the name/ID for the new entity. // DocumentName / ID
	name := "sampletask4"
	// Creates a Key instance.
	taskKey := datastore.NameKey(kind, name, nil)

	// Creates a Task instance.
	doc := Document{
		Status: "CREATED",
		Task: task,
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &doc); err != nil {
		log.Fatalf("Failed to save doc: %v", err)
	}

	fmt.Printf("Saved %v: %v\n", taskKey, doc.Status)
}