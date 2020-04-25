package cloud

import "log"

// client for HTTP CFN endpoint
// NewClient(endpoint,token) *Client
// func (c *Client) Run(task) (Task, err)
// func (c *Client) Delete(task) err
// func (c *Client) GetStatus(task_id) (Task, err)
// func (c *Client) ListTasks() (Tasks, err)

// The client should also be usable as curl alternative for containers to report status/progress
// idea: if nothing was reported back, try to fetch last line of stackdriver output of container -> DataStore.LastMsg

// CFNClient holds infos required to act as client to cloud-vm-docker HTTP cloud function
type CFNClient struct {
	Endpoint    string
	AccessToken string
}

func NewCFNClient(endpoint, accessToken string) *CFNClient {
	return &CFNClient{
		Endpoint:    endpoint,
		AccessToken: accessToken,
	}
}

func (c *CFNClient) Run(taskArgs TaskArguments) (Task, error) {
	log.Printf("CFNClient running on %s with taskArgs %v", c.Endpoint, taskArgs)
	return Task{}, nil
}
