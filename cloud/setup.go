package cloud

func Setup(projectID string) {
	println("Setting up infrastructure for cloud-task-zip-zap ...")
	err := createPubSubTopic(projectID, "ctzz-task-queue")
	if err == nil {
		println("- PubSub Topic created successfully")
	} else {
		println(err.Error())
	}
}

func Destroy(projectID string) {
	println("Destroying infrastructure for cloud-task-zip-zap ...")
	deletePubSubTopic(projectID, "ctzz-task-queue")
}
