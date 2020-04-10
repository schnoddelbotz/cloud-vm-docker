package cloud

func Setup() {
	println("Setting up infrastructure for cloud-task-zip-zap ...")
	err := createPubSubTopic("hacker-playground-254920", "ctzz-task-queue")
	println(err.Error())
}
