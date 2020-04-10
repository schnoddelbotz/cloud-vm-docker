# cloud-task-zip-zap

Probably the easiest way to let dockerized compute tasks run on Google ComputeEngine VMs.
Why cloud-task-zip-zap? Because ...
- CloudFunctions have a maximum runtime of 9 minutes and are limited to 1 or 2 cores and 2 GB of RAM.
  You can only execute code within specific runtimes (Python, Node, Go, ...).
- CloudRun enables running arbitrary Docker images like CloudFunctions, but are also limited (2 CPU/2 GB)
- We just want to run on our dockerized code on ephemeral VMs, without worrying about provisioning
  the VMs nor setting up docker, nor starting processes in those VMs. ctzz does it all!

for deployment, ensure you did this once:
```bash
gcloud auth login
gcloud projects list
gcloud config set project ... 
```

example usage:
- enable datastore in your project via console, then ...
```bash
# set up infrastructure needed (svc acc, cfns, pubsub topics) 
ctzz deploy

# submit a task short running task, wait for and print output!
$ ctzz submit --image busybox --command "echo 'hello world'" --wait --print-logs
2020-01-01 12:30 hello world
$

# submit a long-running task!
$ ctzz submit --image busybox --command "echo 'hello 'world' && sleep 1000 && echo done"
{"uuid":"1234-1234", "status","created"}

# wait for long-running task to complete
$ ctzz wait 1234-1234
Still waiting for 1234-1234 to complete (100s done; no progress feedback, no eta)
Still waiting for 1234-1234 to complete (200s done ...)
...
Task 1234-1234 completed with:
SUCCESS
$ echo $?
0
$

# get logs of given task (fetches VM container logs from stackdriver)
$ ctzz logs 1234-1234

# abort a running task by deleting the VM
$ ctzz kill 1234-1234
```