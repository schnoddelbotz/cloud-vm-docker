package cloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

type Task struct {
	TaskArguments TaskArguments
	Status        string
	VMID          string
	ShutdownToken string
	CreatedAt     time.Time
}

// https://cloud.google.com/datastore/docs/reference/libraries
// https://cloud.google.com/datastore/docs/concepts/queries

// StoreTask saves a new task in FireStore DB.
// It is called by PubSubFn for each message received.
// Note that just storing a task does nothing; the pubsubtrigger is supposed to spin up the VM for the Task.
func StoreTask(projectID string, taskArguments TaskArguments) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name/ID for the new entity. // DocumentName / ID
	name := generateTaskName(taskArguments)

	// Creates a Key instance.
	taskKey := datastore.NameKey(settings.FireStoreCollection, name, nil)
	// Creates a Task instance.
	doc := Task{
		Status:        "CREATED",
		TaskArguments: taskArguments,
		VMID:          "ICH-1000",
		CreatedAt:     time.Now(),
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
	docList := []Task{}
	q := datastore.NewQuery(settings.FireStoreCollection)
	_, err = client.GetAll(ctx, q, &docList)
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	// todo: dynamically check/set required field width
	fmt.Printf("VM_ID          IMAGE                  COMMAND                        CREATED        STATUS\n")
	for _, doc := range(docList) {
		cmd := strings.Join(doc.TaskArguments.Command, " ")
		fmt.Printf("%-14s %-22s %-30s %-14s %s\n", doc.VMID, doc.TaskArguments.Image, cmd, "5 min ago", doc.Status)
	}

	// log.Printf("Got %d tasks as response, showed X, 3 running, 2 deleted.", len(something = client.GetAll retval))
}

func generateTaskName(task TaskArguments) string {
	now := time.Now()
	datePart := now.Format("2006-01-02_15:04:05")
	return fmt.Sprintf("%s_%s", datePart, task.Image)
}