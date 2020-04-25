package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
)

// CloudVMDocker HTTP CloudFunction handler makes VMs triggerable via plain https+token request
func CloudVMDocker(w http.ResponseWriter, r *http.Request, env *Environment) {
	w.Header().Set("Content-Type", "text/plain")
	// TODO:
	// This should become an HTTP entrypoint for exposing cloud-vm-docker functionality via simple JSON api.
	// While using cloud-vm-docker is obviously the most simple/direct approach to submit tasks or
	// manage them -as it speaks to google services like GCE and FireStore directly-
	// it may be helpful to have a RESTish entrypoint for lightweight submission
	// scenarios relying entirely on e.g. just curl.
	// Obviously, it should be (auto-generated-if-not-provided) token-protected, or
	// if the user chooses so, only callable with valid IAM credentials (non-public http endpoint).
	log.Printf("Got request ... %s -> %s", env.GoogleSettings.ProjectID, r.RequestURI)
	//log.Printf("X-Auth-HDR: %s", r.Header.Get("X-Authorization"))
	//log.Printf("REnv: %v", env)

	// hack to see it working. fixme. now. auth-check, everything.
	//  /status/tec1980be9d/BOOTED
	//  /delete/t23c7ac6d4f
	rqParts := strings.Split(r.RequestURI, "/")
	if len(rqParts) == 4 && (rqParts[1] == "delete" || rqParts[1] == "status") {
		action := rqParts[1]
		vmID := rqParts[2]
		targetValue := rqParts[3]

		taskData, err := cloud.GetTask(env.GoogleSettings.ProjectID, vmID)
		if err != nil {
			log.Printf("Error loading task: %s", err)
			http.Error(w, err.Error(), 400)
			return
		}

		if taskData.ManagementToken != r.Header.Get("X-Authorization") {
			log.Printf("DENIED: Invalid token: %s", r.Header.Get("X-Authorization"))
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if action == "delete" {
			err := cloud.DeleteInstanceByName(env.GoogleSettings, vmID)
			if err != nil {
				log.Printf("Error on DeleteInstanceByName(..., %s): %s", vmID, err)
				http.Error(w, err.Error(), 400)
				return
			}
			exitCode, _ := strconv.Atoi(targetValue)
			cloud.UpdateTaskStatus(env.GoogleSettings.ProjectID, vmID, "DONE", exitCode)
		} else if action == "status" {
			cloud.UpdateTaskStatus(env.GoogleSettings.ProjectID, vmID, targetValue)
		} else {
			// todo: add /progress/vmid/99
			log.Printf("Get rid of this if-else shit or I forget myself. Action not implemented.")
		}
		fmt.Fprintf(w, `Thanks for your %s request -- processed successfully`, action)
	} else {
		http.Error(w, "Nope, somehow not.", 400)
	}
}
