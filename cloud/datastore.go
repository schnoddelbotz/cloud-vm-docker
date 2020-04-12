package cloud

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/schnoddelbotz/cloud-task-zip-zap/settings"
)

// https://cloud.google.com/datastore/docs/reference/libraries
// https://cloud.google.com/datastore/docs/concepts/queries

// Thing below should be split into NewTask(proj, args).Store()

// StoreNewTask saves a new task in FireStore DB.
// It is called by PubSubFn for each message received.
// Note that just storing a task does nothing; the pubsubtrigger is supposed to spin up the VM for the Task.
func StoreNewTask(projectID string, taskArguments TaskArguments) Task {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name/ID for the new entity. // DocumentName / ID
	name := generateVMID()
	// Creates a Key instance.
	taskKey := datastore.NameKey(settings.FireStoreCollection, name, nil)
	// Creates a Task instance.
	doc := Task{
		Status:        "CREATED",
		TaskArguments: taskArguments,
		VMID:          name, // dup! also doc title now ...
		CreatedAt:     time.Now(),
		ShutdownToken: generateShutdownToken(),
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &doc); err != nil {
		log.Fatalf("Failed to save doc: %v", err)
	}

	log.Printf("Saved %v: %v\n", taskKey, doc.Status)
	return doc
}

// ListTasks provides 'docker ps' functionality by querying DataStore
func ListTasks(projectID string) {
	// log.Printf("Listing tasks in project %s on collection %s", projectID, settings.FireStoreCollection)
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	docList := make([]Task, 0)
	q := datastore.NewQuery(settings.FireStoreCollection)
	_, err = client.GetAll(ctx, q, &docList)
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	// todo: dynamically check/set required field width -- and/or add flag: disable shortening below...
	fmt.Printf("VM_ID        IMAGE                   COMMAND                                  CREATED        STATUS\n")
	for _, doc := range docList {
		cmd := strings.Join(doc.TaskArguments.Command, " ")
		fmt.Printf("%-12s %-23s %-40s %-14s %s\n",
			doc.VMID, getImageNameWithoutRegistryAndTag(doc.TaskArguments.Image),
			shortenToMaxLength(cmd, 38), "5 min ago", doc.Status)
	}

	// log.Printf("Got %d tasks as response, showed X, 3 running, 2 deleted.", len(something = client.GetAll retval))
}

func generateTaskName(task TaskArguments) string {
	// ... WAS: as used as DataStore ID -- drop?
	now := time.Now()
	datePart := now.Format("2006-01-02_15:04:05")
	return fmt.Sprintf("%s_%s", datePart, task.Image)
}

func generateVMID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("t%x", b[0:5])
}

func generateShutdownToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func getImageNameWithoutRegistryAndTag(imageFQDN string) string {
	imageParts := strings.Split(imageFQDN, "/")
	imageNameWithTag := imageParts[len(imageParts)-1]
	return strings.Split(imageNameWithTag, ":")[0]
}

func shortenToMaxLength(str string, max int) string {
	// shorten String to maximum length. If cut, precede with ellipsis.
	from := len(str) - max
	if from < 0 {
		from = 0
	}
	result := str[from:len(str)]
	if len(str) > max && max > 3 {
		result = "…" + result[1:]
	}
	return result
}
