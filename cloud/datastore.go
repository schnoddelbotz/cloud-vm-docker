package cloud

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// Thing below should be split into NewTask(proj, args).Store()

// StoreNewTask saves a new task in FireStore DB.
// It is called by PubSubFn for each message received.
// Note that just storing a task does nothing; the pubsubtrigger is supposed to spin up the VM for the Task.
func StoreNewTask(projectID string, taskArguments TaskArguments) Task {
	ctx := context.Background()
	client := NewDataStoreClient(ctx, projectID)
	// Sets the name/ID for the new entity. // DocumentName / ID
	taskKey := datastore.NameKey(settings.FireStoreCollection, taskArguments.VMID, nil)
	doc := Task{
		Status:          TaskStatusCreated,
		TaskArguments:   taskArguments,
		VMID:            taskArguments.VMID, // dup! also doc title now ...
		CreatedAt:       time.Now(),
		ManagementToken: generateManagementToken(),
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
	ctx := context.Background() // fixme pass in
	client := NewDataStoreClient(ctx, projectID)
	docList := make([]Task, 0)

	// todo: add filter (-a arg), sorting, FIX CREATED OUTPUT
	q := datastore.NewQuery(settings.FireStoreCollection)
	_, err := client.GetAll(ctx, q, &docList)
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	// todo: dynamically check/set required field width -- and/or add flag: disable shortening below...
	outputFormat := "%-12s %-23s %-40s %-14s %s\n"
	fmt.Printf(outputFormat, "VM_ID", "IMAGE", "COMMAND", "CREATED", "STATUS")
	for _, doc := range docList {
		cmd := strings.Join(doc.TaskArguments.Command, " ")
		fmt.Printf(outputFormat,
			doc.VMID, getImageNameWithoutRegistryAndTag(doc.TaskArguments.Image),
			shortenToMaxLength(cmd, 38), "5 min ago", doc.Status)
	}

	// log.Printf("Got %d tasks as response, showed X, 3 running, 2 deleted.", len(something = client.GetAll retval))
}

// NewDataStoreClient returns a dataStore client and its context, exits fatally on error
func NewDataStoreClient(ctx context.Context, projectID string) *datastore.Client {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

// UpdateTaskStatus sets a Task's 'status' field to given string value
func UpdateTaskStatus(projectID, vmID, status string) error {
	log.Printf("UpdateTaskStatus: %s -> %s", vmID, status)
	ctx := context.Background() // fixme pass in
	client := NewDataStoreClient(ctx, projectID)
	tx, err := client.NewTransaction(ctx)
	if err != nil {
		log.Fatalf("client.NewTransaction: %v", err)
	}
	var task Task
	taskKey := datastore.NameKey(settings.FireStoreCollection, vmID, nil)
	if err := tx.Get(taskKey, &task); err != nil {
		log.Fatalf("tx.Get: %v", err)
	}
	task.Status = status
	if _, err := tx.Put(taskKey, &task); err != nil {
		log.Fatalf("tx.Put: %v", err)
	}
	if _, err := tx.Commit(); err != nil {
		log.Fatalf("tx.Commit: %v", err)
	}
	return nil
}

// GetTask tries to fetch given record from DataStore
func GetTask(projectID, vmID string) (Task, error) {
	ctx := context.Background() // fixme pass in
	client := NewDataStoreClient(ctx, projectID)
	var task Task
	taskKey := datastore.NameKey(settings.FireStoreCollection, vmID, nil)
	err := client.Get(ctx, taskKey, &task)
	return task, err
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

func generateManagementToken() string {
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
