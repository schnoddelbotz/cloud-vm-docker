package api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/schnoddelbotz/cloud-vm-docker/cloud"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// client for HTTP CFN endpoint
// TODO:
// func (c *Client) Delete(task) err
// func (c *Client) GetStatus(task_id) (Task, err)
// func (c *Client) ListTasks() (Tasks, err)

// The client should also be usable as curl alternative for containers to report status/progress
// idea: if nothing was reported back, try to fetch last line of stackdriver output of container -> FireStore.LastMsg

// CFNClient holds infos required to act as client to cloud-vm-docker HTTP cloud function
type CFNClient struct {
	Endpoint    string
	AccessToken string
	HTTPClient  *http.Client
}

// NewCFNClient gives a new client for followup CloudFunction requests
func NewCFNClient(endpoint, accessToken string) *CFNClient {
	client := &CFNClient{
		Endpoint:    endpoint,
		AccessToken: accessToken,
	}
	client.HTTPClient = &http.Client{
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	return http.ErrUseLastResponse
		//},
	}
	return client
}

// GetEndpoint returns the HTTP URL of deployed CloudVMDocker CloudFunction, based on project and region
func GetEndpoint(project, region string) string {
	if project == "" {
		log.Fatalf("Cannot create CFN HTTP Client without --project")
	}
	if region == "" {
		log.Fatalf("Cannot create CFN HTTP Client without --region")
	}
	return fmt.Sprintf("https://%s-%s.cloudfunctions.net/CloudVMDocker", region, project)
}

// Run submits a Docker task to run in the cloud to HTTP CloudFunction endpoint
func (c *CFNClient) Run(taskArgs cloud.TaskArguments) (cloud.Task, error) {
	log.Printf("CFNClient running on %s with taskArgs %v", c.Endpoint, taskArgs)
	requestBody, err := json.Marshal(taskArgs)
	if err != nil {
		log.Fatalf("Error building requests JSON: %s", err)
	}
	requestPath := fmt.Sprintf("run/%s/%s", taskArgs.VMID, taskArgs.VMType)
	responseBody, err := c.executeClientRequest("POST", requestPath, requestBody)
	if err != nil {
		log.Fatalf("Error submitting request: %s", err)
	}
	var task cloud.Task
	err = json.Unmarshal(responseBody, &task)
	if err != nil {
		log.Printf("Unmarshalling response failed: %s", err)
	}
	log.Printf("SUCCESS: Got and returning task: %v", task)
	return task, nil
}

func (c *CFNClient) WaitForDoneStatus(vmID string) cloud.Task {
	var task cloud.Task
	status := "unknown"
	requestPath := fmt.Sprintf("status/%s/", vmID)

	log.Printf("Waiting for DONE status of VM %s", vmID)
	for status != "DONE" {
		body, err := c.executeClientRequest("GET", requestPath, []byte{})
		if err != nil {
			log.Fatalf("ClientERR: %s", err)
		}
		err = json.Unmarshal(body, &task)
		if err != nil {
			log.Fatalf("Client JSON ERR: %s", err)
		}
		if status != task.Status {
			log.Printf("Status change: %s -> %s", status, task.Status)
		}
		status = task.Status
		if status != "DONE" {
			time.Sleep(30 * time.Second)
		}
	}
	log.Printf("SUCCESS waiting for DONE status of VM %s", vmID)
	return task
}

func (c *CFNClient) executeClientRequest(method, path string, requestBody []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.Endpoint, path)
	requestBodyReader := bytes.NewReader(requestBody)
	request, err := http.NewRequest(method, url, requestBodyReader)
	if err != nil {
		log.Fatalf("Client failed: %s", err)
		return []byte{}, err
	}
	request.Header.Set("X-Authorization", c.AccessToken)
	request.Header.Set("Content-type", "application/json")

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		log.Fatalf("Client failed: %s", err)
	}
	if response.StatusCode != 200 {
		log.Printf("Client got non-200 response: %d", response.StatusCode)
		log.Printf("LOC: %s", response.Header)
		// 302: redirect by google to auth page, if CFN is not deployed or not public
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	return responseBody, err
}
