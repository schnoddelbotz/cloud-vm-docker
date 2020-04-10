package cloud

func Setup() {
	println("Setting up infrastructure for cloud-task-zip-zap ...")
	createPubSubTopic("ctzz-task-queue")
}
