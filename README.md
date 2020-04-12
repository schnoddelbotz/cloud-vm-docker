# cloud-task-zip-zap

Running dockerized one-shot workloads on Google ComputeEngine has never been easier.
At least this might do as a valid advertisement slogan for cloud-task-zip-zap, given:

```bash
# run a task locally, using plain, local Docker
$ docker run busybox echo foo 

# the same, but run task on Docker on a ComputeEngine VM
$ ctzz run busybox echo foo
```

## usage

OK. That looked too simple, as it was using all the defaults from environment.
So, a more complex example:

```bash
# run command from custom GCR-hosted image, using a VM with 16 cores
$ ctzz run \
    -d \
    --vm-type n1-standard-16 \
    eu.gcr.io/my-project-6afd9bfb/my-compute-task-image:latest \
    bash -c "echo started && sleep 3600 && echo done"

# Like `docker run -d`, the above command will not wait for task to complete
# and will not print logs. Of course, they're accessible +/- as if it was plain Docker.
$ ctzz ps
VM_ID        IMAGE               COMMAND                  CREATED             STATUS
6af7db3a     eu.gcr.io.my-pro... sleep 3600               15 min ago          running

# Containers running on VMs will forward logs to StackDriver. To read those logs, like in Docker, do:
$ ctzz logs 6af7db3a
2020/04/12 10:20:05 started

# Compute tasks are best run in forground (e.g. in Airflow DAGs), as this will implicitly wait
# for container command completion. 
# But if you decided  to run a task 'detached' (-d), then you can wait for completion:
$ ctzz wait 6af7db3a
2020/04/12 10:30:15 started waiting for completion of task 6af7db3a
2020/04/12 10:35:10 task 6af7db3a completed - setting wait's exit status to the task's one: EXIT_STATUS_OK
```

## why?

Why bother with cloud-task-zip-zap? Because ...
- CloudFunctions have a maximum runtime of 9 minutes and are limited to 1 or 2 cores and 2 GB of RAM.
  You can only execute code within specific runtimes (Python, Node, Go, ...).
- CloudRun enables running arbitrary Docker images in CloudFunctions-style, but are also limited (2 CPU/2 GB)
- We just want to run on our Dockerized code on ephemeral VMs, without neither worrying about provisioning
  the VMs nor setting up Docker, nor starting processes in those VMs -- ctzz does it all!

## setting up google cloud for ctzz usage

for deployment, ensure you did this once:
```bash
gcloud auth login
gcloud projects list
gcloud config set project ... 

# to let run ctzz locally and interact with google services, create a svc account as in
# https://cloud.google.com/datastore/docs/reference/libraries

# best in your .bashrc
export GOOGLE_APPLICATION_CREDENTIALS=$HOME/.config/gcloud/svc-account.json

# one day, `ctzz setup` should do, but for now ... rely on gcloud. could use docker image...
make gcp_deploy

# the above command will deploy TWO cloud functions, including PubSub Topic and subscription:
# - one HTTP endpoint, intended for "dumb" curl clients to submit VM tasks (by writing PubSub message)
# - one PubSub monitor, that will spin up VMs as requested in PubSub messages
```

## TODO

- have some monitoring dashboard web endpoint using `status` data + google monitoring/logs links ... 
- or update some google-hosted dashboard to add/remove machines as they come/run/go(history)
- auto-completion
- let user decide on `run` whether to use HTTP endpoint or write pubsub directly
- deployment: let user disable HTTP endpoint if not needed
- list which commands work as 100% "drop-in" replacement for docker commands -- goal: as-much-as-possible