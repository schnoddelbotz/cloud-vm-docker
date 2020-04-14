# cloud-vm-docker -- !!WIP!!

Running dockerized one-shot workloads on Google ComputeEngine has never been easier.
At least this might do as a valid advertisement slogan for cloud-vm-docker, given:

```bash
# run a task locally, using plain, local Docker
$ docker run busybox echo foo

# the same, but run task on Docker on a ComputeEngine VM
$ cloud-vm-docker run busybox echo foo
```

## (intended) usage

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

### example use cases

- use cloud-vm-docker eg. in your Airflow workflows, to off-load resource hungry compute tasks to the cloud
- use cloud-vm-docker eg. to run some jMeter benchmarks on capable cloud VMs ... against your own site
- use cloud-vm-docker eg. data-heavy processing tasks, which benefit from cloud data "locality"
- use cloud-vm-docker eg. to spin up an VM instance for further operations to be carried out via SSH
- playing with Go, Docker and GoogleCloud APIs

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

## test-drive -- what works now?

Git clone this repo, and adjust `testenv.inc.sh` to your needs. Then ...

```bash
# especially adjust command for getting auto-completion :-)
source testenv.inc.sh
make deploy_gcp
make clean test build

# this creates a VM directly (from local machine), bypassing PubSub
./cloud-vm-docker task-vm create busybox sh -c 'echo hello world ; sleep 120 ; echo goodnight'

# ^^ notice:
# - ComputeEngine console UI should show the VM within a few secs
# - If nothing goes wrong (TM), the VM should self-destruct upon completion, just leaving logs

# the same, but using "official" way via PubSub Message (which is processed by a CFn)
# NOTE: Does NOT spawn the VM atm, just logs what it will do soon...
./cloud-vm-docker run busybox sh -c 'echo hello world ; sleep 120 ; echo goodnight'

# list VMs as stored in dataStore
./cloud-vm-docker ps

# this should be ./cloud-vm-docker ssh ... but, for now, look up IP in console[FIXME].
# If this fails ... there's a bug with more than 1 ssh keys in home. +1 fixme...
ssh cloud-vm-docker@IP 

# delete the VM -- [TODO: autocomplete!]
./cloud-vm-docker task-vm kill ...ID_as_shown_in_ps_output...
``` 

## links

Google Cloud general

- https://cloud.google.com/compute/docs/regions-zones#available

VMs

- https://github.com/googleapis/google-api-go-client/blob/master/examples/compute.go
- https://cloud.google.com/compute/docs/reference/rest/v1/instances/insert
- https://godoc.org/google.golang.org/api/compute/v1
- https://cloudinit.readthedocs.io/en/latest/index.html
- https://www.freedesktop.org/software/systemd/man/systemd.service.html

DataStore

- https://cloud.google.com/datastore/docs/reference/libraries
- https://cloud.google.com/datastore/docs/concepts/queries

PubSub

- https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/pubsub/pubsub.go

Operations

- https://github.com/googleapis/google-api-go-client/blob/master/examples/operation_progress.go

## TODO

- tests, tests, tests
- DataStore: updates -- status updates via CFN (see cloud_init), or after `task-vm kill ...`
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