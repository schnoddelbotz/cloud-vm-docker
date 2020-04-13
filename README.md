# cloud-vm-docker

Running dockerized one-shot workloads on Google ComputeEngine has never been easier.
At least this might do as a valid advertisement slogan for cloud-vm-docker, given:

```bash
# run a task locally, using plain, local Docker
$ docker run busybox echo foo

# the same, but run task on Docker on a ComputeEngine VM
$ cloud-vm-docker run busybox echo foo
```

## usage

OK. That looked too simple, as it was using all the defaults from environment.
So, a more complex example:

```bash
# run command from custom GCR-hosted image, using a VM with 16 cores
$ cloud-vm-docker run \
    -d \
    --vm-type n1-standard-16 \
    eu.gcr.io/my-project-6afd9bfb/my-compute-task-image:latest \
    bash -c "echo started && sleep 3600 && echo done"

# Like `docker run -d`, the above command will not wait for task to complete
# and will not print logs. Of course, they're accessible +/- as if it was plain Docker.
$ cloud-vm-docker ps
VM_ID       IMAGE                   COMMAND                                  CREATED        STATUS
fb0f979473  busybox                 echo foo                                 5 min ago      created

# Containers running on VMs will forward logs to StackDriver. To read those logs, like in Docker, do:
$ cloud-vm-docker logs fb0f979473
2020/04/12 10:20:05 started

# Compute tasks are best run in forground (e.g. in Airflow DAGs), as this will implicitly wait
# for container command completion.
# But if you decided  to run a task 'detached' (-d), then you can wait for completion:
$ cloud-vm-docker wait fb0f979473
2020/04/12 10:30:15 started waiting for completion of task 6af7db3a
2020/04/12 10:35:10 task 6af7db3a completed - setting wait's exit status to the task's one: EXIT_STATUS_OK
```

Note that, unless disabled by using `--no-ssh`, cloud-vm-docker will automatically "upload"
(using cloud_init) your SSH public keys from `~/.ssh/*.pub`, which will authorize you for SSH VM logins.
To come: `cloud-vm-docker ssh <my-vm> [<CMD>]`.

## why?

Why bother with cloud-vm-docker? Because ...
- CloudFunctions have a maximum runtime of 9 minutes and are limited to 1 or 2 cores and 2 GB of RAM.
  You can only execute code within specific runtimes (Python, Node, Go, ...).
- CloudRun enables running arbitrary Docker images in CloudFunctions-style, but are also limited (2 CPU/2 GB)
- We just want to run on our Dockerized code on ephemeral VMs, without neither worrying about provisioning
  the VMs nor setting up Docker, nor starting processes in those VMs -- cloud-vm-docker does it all!
- To circumvent some limitations listed in https://cloud.google.com/compute/docs/containers/deploying-containers

## setting up google cloud for cloud-vm-docker usage

for deployment, ensure you did this once:
```bash
gcloud auth login
gcloud projects list
gcloud config set project ...

# to let run cloud-vm-docker locally and interact with google services, create a svc account as in
# https://cloud.google.com/datastore/docs/reference/libraries

# best in your .bashrc
export GOOGLE_APPLICATION_CREDENTIALS=$HOME/.config/gcloud/svc-account.json

# one day, `cloud-vm-docker setup` should do, but for now ... rely on gcloud. could use docker image...
make gcp_deploy

# the above command will deploy TWO cloud functions, including PubSub Topic and subscription:
# - one HTTP endpoint, intended for "dumb" curl clients to submit VM tasks (by writing PubSub message)
# - one PubSub monitor, that will spin up VMs as requested in PubSub messages
```

## links

- https://cloud.google.com/compute/docs/regions-zones#available

VMs

- https://github.com/googleapis/google-api-go-client/blob/master/examples/compute.go
- https://cloud.google.com/compute/docs/reference/rest/v1/instances/insert
- https://godoc.org/google.golang.org/api/compute/v1
- https://cloudinit.readthedocs.io/en/latest/index.html
- https://www.freedesktop.org/software/systemd/man/systemd.service.html

Operations

- https://github.com/googleapis/google-api-go-client/blob/master/examples/operation_progress.go

## TODO

- tests, tests, tests
- have some monitoring dashboard web endpoint using `status` data + google monitoring/logs links ...
- or update some google-hosted dashboard to add/remove machines as they come/run/go(history)
- https://cloud.google.com/compute/docs/storing-retrieving-metadata --> put VM meta in DataStore / partially?
  ```bash
  curl -H'Metadata-Flavor:Google' "http://metadata.google.internal/computeMetadata/v1/instance/"curl -H'Metadata-Flavor:Google' "http://metadata.google.internal/computeMetadata/v1/instance/"
  curl -H'Metadata-Flavor:Google' "http://metadata.google.internal/computeMetadata/v1/instance/attributes/user-data"
  ```
- let user decide on `run` whether to use HTTP endpoint or write pubsub directly
- deployment: let user disable HTTP endpoint if not needed
- list which commands work as 100% "drop-in" replacement for docker commands -- goal: as-much-as-possible
- coool! can I use this for interactive containers as well? no, not yet, maybe never. you can ssh to vm though.
- allow alternate VM disk images? custom cloud_init? custom network? labels? svcAccount (or roles to add to default)?
- have some simple dashboard ('docker ps++') served via http cfn?