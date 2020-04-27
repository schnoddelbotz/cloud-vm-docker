package cloud

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

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
	tasksRef, err := client.Collection(settings.FireStoreCollection).
		OrderBy("CreatedAt", firestore.Asc).
		//Limit(25). // FIXME: no-constant-magic ... make user flag
		Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	for _, i := range tasksRef {
		t, err := fireStoreDataToTask(i.Data())
		if err != nil {
			log.Printf("FireStore data (%v) error: %s", t, err)
			continue
		}
		docList = append(docList, t)
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
	task, err = fireStoreDataToTask(d.Data())
	if err != nil {
		log.Fatalf("Ooops. Cannot convert FireStore data to task %s: %s", vmID, err)
	}
	return task, err
}

func fireStoreDataToTask(dict map[string]interface{}) (Task, error) {
	// this feels wrong from the start. fixme.
	var task Task
	jsonbody, err := json.Marshal(dict)
	if err != nil {
		return task, err
	}
	err = json.Unmarshal(jsonbody, &task)
	return task, err
}

func generateTaskName(task TaskArguments) string {
	// ... WAS: as used as FireStore ID -- drop?
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
