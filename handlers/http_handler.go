package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
)

// CloudVMDocker HTTP CloudFunction handler makes VMs triggerable via plain https+token request
func CloudVMDocker(w http.ResponseWriter, r *http.Request, env *Environment) {
	// TODO: Support both tokens for VM operations: per-VM one and "general" one (for CFN)
	//  MGM-TOKEN POST /run/tec1980be9d/fooBarBaz
	//   VM-TOKEN GET /status/tec1980be9d/BOOTED --- semantically wrong to use GET, but easier to curl... m(
	//  MGM-TOKEN GET /status/tec1980be9d
	//  MGM-TOKEN GET /status (cloud-vm-docker ps)
	//   VM-TOKEN GET /delete/t23c7ac6d4f/0  --- semantically wrong to use GET, but easier to curl... m(
	// FIXMEEEEE! Separete /vm (vm management token) and /manage (cfn admin token)
	action, vmID, targetValue, err := parseRequestURI(r.RequestURI)
	if err != nil {
		log.Printf("Invalid request URI: %s", r.RequestURI)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientToken := r.Header.Get("X-Authorization")

	switch action {
	case "status":
		// todo: distinguish different status requests (paths) ...
		if r.Method == http.MethodPost {
			if !authenticateVMRequest(w, clientToken, env, vmID) {
				return
			}
			cloud.UpdateTaskStatus(env.GoogleSettings.ProjectID, vmID, targetValue)
		} else {
			if !authenticateAdminRequest(w, clientToken, env) {
				return
			}
			handleStatusGet(w, env, vmID)
		}
	case "delete":
		if !authenticateVMRequest(w, clientToken, env, vmID) {
			return
		}
		handleDelete(w, env, vmID, targetValue)
	case "container":
		if !authenticateVMRequest(w, clientToken, env, vmID) {
			return
		}
		handleContainer(w, r, env, vmID)
	case "run":
		if !authenticateAdminRequest(w, clientToken, env) {
			return
		}
		handleRun(w, r, env, vmID)
	}
}

func authenticateAdminRequest(w http.ResponseWriter, clientToken string, env *Environment) bool {
	if clientToken != env.GoogleSettings.AccessToken {
		log.Printf("Permission denied for admin request - bad token: %s", clientToken)
		http.Error(w, "FIXME This should be a JSON 401", 401)
		return false
	}
	return true
}

func authenticateVMRequest(w http.ResponseWriter, clientToken string, env *Environment, vmID string) bool {
	taskData, err := cloud.GetTask(env.GoogleSettings.ProjectID, vmID)
	if err != nil {
		log.Printf("Error loading task: %s", err)
		http.Error(w, err.Error(), 500)
		return false
	}

	if taskData.ManagementToken != clientToken {
		log.Printf("DENIED: Invalid token: %s", clientToken)
		http.Error(w, err.Error(), http.StatusForbidden)
		return false
	}
	return true
}

func parseRequestURI(uri string) (action, vmid, data string, err error) {
	parts := strings.Split(uri, "/")
	if len(parts) != 4 {
		err = errors.New("invalid URI, expected /CloudVMDocker/ACTION/VM_ID/DATA")
		return
	}
	action = parts[1]
	if action != "delete" && action != "status" && action != "run" && action != "container" {
		err = errors.New("invalid action %s, expected on of: delete, status, run, container")
		return
	}
	vmid = parts[2]
	// todo: validate VMID: starts with t, len=...12?
	data = parts[3]
	return
}

func handleRun(w http.ResponseWriter, r *http.Request, env *Environment, vmID string) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error loading request body: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	taskArguments := cloud.NewTaskArgumentsFromBytes(requestBody)
	log.Printf("Writing task to DataStore: %+v", taskArguments)
	task := cloud.StoreNewTask(env.GoogleSettings.ProjectID, *taskArguments)
	createOp, err := cloud.CreateVM(env.GoogleSettings, task)
	if err != nil {
		log.Printf("ERROR running TaskArguments: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println("VM creation requested successfully, waiting for op...")
	cloud.WaitForOperation(env.GoogleSettings.ProjectID, env.GoogleSettings.Zone, createOp.Name)

	log.Printf("Saving GCE InstanceID to DataStore: %s => %d", vmID, createOp.TargetId)
	err = cloud.SetTaskInstanceId(env.GoogleSettings.ProjectID, vmID, createOp.TargetId)
	if err != nil {
		log.Printf("ARGH!!! Could not update instanceID in DataStore: %s", err)
	}
	task.InstanceID = strconv.FormatUint(createOp.TargetId, 10)

	responseBody, _ := json.Marshal(task)
	w.Header().Set("content-type", "application/json")
	numBytes, err := w.Write(responseBody)
	if err != nil {
		log.Fatalf("Waaaah! Failed sending off %d bytes to client, who will be unhappy for sure: %s", numBytes, err)
	}
}

func handleStatusGet(w http.ResponseWriter, env *Environment, vmID string) {
	log.Printf("Serving status for vm %s", vmID)
	task, err := cloud.GetTask(env.GoogleSettings.ProjectID, vmID)
	if err != nil {
		log.Printf("Error handling task status get requests vmid=%s: %s", vmID, err)
		http.Error(w, err.Error(), 500)
		return
	}
	repsonseBody, _ := json.Marshal(task)
	w.Header().Set("content-type", "application/json")
	w.Write(repsonseBody)
}

func handleDelete(w http.ResponseWriter, env *Environment, vmID string, exitCodeString string) {
	log.Printf("Handling DELETE request form VMID %s with exitCode %s", vmID, exitCodeString)
	err := cloud.DeleteInstanceByName(env.GoogleSettings, vmID)
	if err != nil {
		log.Printf("Error on DeleteInstanceByName(..., %s): %s", vmID, err)
		http.Error(w, err.Error(), 500)
		return
	}
	exitCode, err := strconv.Atoi(exitCodeString)
	if err != nil {
		log.Printf("Error on DeleteInstanceByName(..., %s): Cannot convert exit code '%s': %s", vmID, exitCodeString, err)
		http.Error(w, err.Error(), 500)
		return
	}
	err = cloud.UpdateTaskStatus(env.GoogleSettings.ProjectID, vmID, "DONE", exitCode)
	if err != nil {
		log.Printf("Error on DeleteInstanceByName(..., %s): Unable to update DataStore after successful VM deletion %s", vmID, err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, `Thanks for your DELETE request -- processed successfully`)
}

func handleContainer(w http.ResponseWriter, r *http.Request, env *Environment, vmID string) {
	log.Printf("Handling CONTAINER(ID) request form VMID %s", vmID)
	containerID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Cannot read body of container update request: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if len(containerID) != 64 {
		log.Printf("ERROR: Bad container update request, expected 64-char container ID, got: `%s`", containerID)
		http.Error(w, err.Error(), 500)
		return
	}
	err = cloud.SetTaskContainerID(env.GoogleSettings.ProjectID, vmID, string(containerID))
	if err != nil {
		log.Printf("ERROR: Failed to update containerID of vm %s to `%s`", vmID, containerID)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, `Thanks for your CONTAINER request -- processed successfully .... not yet..`)
}
