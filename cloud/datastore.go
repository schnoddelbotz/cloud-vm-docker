package cloud

// https://cloud.google.com/datastore/docs/reference/libraries

import (
	"context"
	"fmt"
	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

type Document struct {
	Task Task
	Status string
	VMID string
	ShutdownToken string
}

func StoreTask(projectID string, task Task) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name/ID for the new entity. // DocumentName / ID
	name := generateTaskName(task)

	// Creates a Key instance.
	taskKey := datastore.NameKey(settings.FireStoreCollection, name, nil)
	// Creates a Document instance.
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

func ListTasks(projectID string) {
	// log.Printf("Listing tasks in project %s on collection %s", projectID, settings.FireStoreCollection)
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	docList := []Document{}
	q := datastore.NewQuery(settings.FireStoreCollection)
	_, err = client.GetAll(ctx, q, &docList)
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	// todo: dynamically check/set required field width
	fmt.Printf("VM_ID          IMAGE                  COMMAND                        CREATED        STATUS\n")
	for _, doc := range(docList) {
		cmd := strings.Join(doc.Task.Command, " ")
		fmt.Printf("%-14s %-22s %-30s %-14s %s\n", doc.VMID, doc.Task.Image, cmd, "5 min ago", doc.Status)
	}

	// log.Printf("Got %d tasks as response, showed X, 3 running, 2 deleted.", len(something = client.GetAll retval))
}

func generateTaskName(task Task) string {
	now := time.Now()
	datePart := now.Format("2006-01-02_15:04:05")
	return fmt.Sprintf("%s_%s", datePart, task.Image)
}