package cloud

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/schnoddelbotz/cloud-vm-docker/settings"
)

// Thing below should be split into NewTask(proj, args).Store()

// StoreNewTask saves a new task in FireStore DB.
func StoreNewTask(projectID string, taskArguments TaskArguments) Task {
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	doc := Task{
		Status:          TaskStatusCreated,
		TaskArguments:   taskArguments,
		VMID:            taskArguments.VMID, // dup! also doc title now ...
		CreatedAt:       time.Now(),
		DockerExitCode:  -1,
		ManagementToken: generateManagementToken(),
	}
	// Saves the new entity.
	if _, err := client.Collection(settings.FireStoreCollection).Doc(taskArguments.VMID).Set(ctx, doc); err != nil {
		log.Fatalf("Failed to save doc: %v", err)
	}
	log.Printf("Saved %v: %v", taskArguments.VMID, doc.Status)
	return doc
}

// ListTasks provides 'docker ps' functionality by querying FireStore
func ListTasks(projectID string) {
	ctx := context.Background() // fixme pass in
	client := NewFireStoreClient(ctx, projectID)
	docList := make([]Task, 0)

	// todo: add filter (-a arg), sorting, FIX CREATED OUTPUT
	iter := client.Collection(settings.FireStoreCollection).
		OrderBy("CreatedAt", firestore.Asc).
		//Limit(25). // FIXME: no-constant-magic ... make user flag
		Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("FireStore iterator boom: %s", err)
		}

		var task Task
		if err := doc.DataTo(&task); err != nil {
			log.Printf("FireStore data error: %s", err)
			continue
		}
		docList = append(docList, task)
	}

	// todo: dynamically check/set required field width -- and/or add flag: disable shortening below...
	// todo: add --long option to  ...include console log links?  ...include live gce vm state?
	outputFormat := "%-12s %-23s %-40s %-14s %s\n"
	fmt.Printf(outputFormat, "VM_ID", "IMAGE", "COMMAND", "CREATED", "STATUS")
	for _, doc := range docList {
		cmd := strings.Join(doc.TaskArguments.Command, " ")
		fmt.Printf(outputFormat,
			doc.VMID, getImageNameWithoutRegistryAndTag(doc.TaskArguments.Image),
			shortenToMaxLength(cmd, 38),
			fmt.Sprintf("%.f mins ago", time.Since(doc.CreatedAt).Minutes()),
			fmt.Sprintf("%s (%d)", doc.Status, doc.DockerExitCode))
	}

	// log.Printf("Got %d tasks as response, showed X, 3 running, 2 deleted.", len(something = client.GetAll retval))
}

// NewFireStoreClient returns a dataStore client and its context, exits fatally on error
func NewFireStoreClient(ctx context.Context, projectID string) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

// UpdateTaskStatus sets a Task's 'status' field to given string value
func UpdateTaskStatus(projectID, vmID, status string, exitCode ...int) error {
	log.Printf("UpdateTaskStatus: %s -> %s", vmID, status)
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	update := map[string]interface{}{"Status": status}
	if len(exitCode) > 0 {
		update["DockerExitCode"] = exitCode[0]
	}
	_, err := client.Collection(settings.FireStoreCollection).Doc(vmID).
		Set(ctx, update, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

// SetTaskContainerID sets a Task's 'status' field to given string value
func SetTaskContainerID(projectID, vmID, containerID string) error {
	log.Printf("SetTaskContainerID: %s -> %s", vmID, containerID)
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	_, err := client.Collection(settings.FireStoreCollection).Doc(vmID).
		Set(ctx, map[string]interface{}{"DockerContainerId": containerID}, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTaskStatus sets a Task's 'status' field to given string value
func SetTaskInstanceId(projectID, vmID string, instanceID uint64) error {
	log.Printf("SetTaskInstanceId: VM_ID %s -> InstanceID %d", vmID, instanceID)
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	_, err := client.Collection(settings.FireStoreCollection).Doc(vmID).
		Set(ctx, map[string]interface{}{"InstanceID": strconv.FormatUint(instanceID, 10)}, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

// GetTask tries to fetch given record from FireStore
func GetTask(projectID, vmID string) (Task, error) {
	ctx := context.Background() // fixme pass in
	client := NewFireStoreClient(ctx, projectID)
	var task Task
	d, err := client.Collection(settings.FireStoreCollection).Doc(vmID).Get(ctx)
	if err != nil {
		return Task{}, err
	}
	if err := d.DataTo(&task); err != nil {
		log.Fatalf("Ooops. Cannot convert FireStore data to task %s: %s", vmID, err)
	}
	return task, err
}

// WaitForTaskDone should use FireStore realtime updates to be notified on updates ...
func WaitForTaskDone(projectID, vmID string) (Task, error) {
	task, err := GetTask(projectID, vmID)
	if err != nil {
		return task, err
	}
	if task.Status == "DONE" {
		log.Printf("Found DONE status for task on initial FireStore load request")
		return task, nil
	}

	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	col := client.Collection(settings.FireStoreCollection) //.Doc(vmID)
	log.Printf("Waiting for task status DONE for vmID %s which is now in status %s", vmID, task.Status)
	//firestore.LogWatchStreams = true
	// why is this all private...???
	// https://github.com/googleapis/google-cloud-go/blob/master/firestore/watch.go
	// https://stackoverflow.com/questions/51200460/how-to-listen-to-firestore-through-rpc -- THANKS, @kataras
	iter := col.Snapshots(ctx)
	defer iter.Stop()
	keepGoing := true
	for keepGoing {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return task, err
		}

		for _, change := range doc.Changes {
			switch change.Kind {
			//case firestore.DocumentRemoved:
			//case firestore.DocumentAdded:
			// isNew := change.Doc.CreateTime.After(task.CreatedAt)
			case firestore.DocumentModified:
				err := change.Doc.DataTo(&task)
				if err != nil {
					log.Fatalf("Cannot read task data: %e", err)
				}
				if task.VMID == vmID {
					log.Printf("Change received. New status: %s", task.Status)
					if task.Status == "DONE" {
						log.Printf("Done waiting. Task DockerExitCode: %d", task.DockerExitCode)
						keepGoing = false
					}
				} else {
					log.Printf("Ignoring change on VMID %s (which is not me! I'm %s)", task.VMID, vmID)
				}
			}
		}
		// log.Printf("NEXT PLEASE") // -- will be passed once for every notification received
	}
	if task.DockerExitCode != 0 {
		return task, fmt.Errorf("task returned with non-zero DockerExitCode: %d", task.DockerExitCode)
	}
	return task, nil
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
