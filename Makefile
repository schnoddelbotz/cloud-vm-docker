
BINARY := cloud-vm-docker

VERSION := $(shell git describe --tags | cut -dv -f2)
DOCKER_IMAGE := schnoddelbotz/cloud-vm-docker
LDFLAGS := -X github.com/schnoddelbotz/cloud-vm-docker/cmd.AppVersion=$(VERSION) -w

GO_SOURCES := */*.go */*/*.go

CVD_TOKEN ?= FIXME-REPLACE-ME-DURING-DEPLOYMENT

build: $(BINARY)

$(BINARY): $(GO_SOURCES)
	# building cloud-vm-docker
	go build -v -o $(BINARY) -ldflags='-w -s $(LDFLAGS)' ./cli/main.go

all_local: clean test build

all_docker: clean test docker_image_prod

release: all_docker docker_image_push


test:
	# golint -set_exit_status ./...
	go fmt ./...
	# golint ./...
	go vet ./...
	go test -ldflags='-w -s $(LDFLAGS)' ./...
	go build cloudfunction.go

coverage: clean
	PROVIDER=MEMORY go test -coverprofile=coverage.out -coverpkg=./... -ldflags='-w -s $(LDFLAGS)' ./...
	go tool cover -html=coverage.out

deploy_gcp: test clean
	gcloud functions deploy CloudVMDocker --region=europe-west1 --runtime go113 \
 		--trigger-http --allow-unauthenticated --project=$(CVD_PROJECT) \
 		--set-env-vars=CVD_PROJECT=$(CVD_PROJECT),CVD_TOKEN=$(CVD_TOKEN),CVD_VERSION=$(VERSION)

docker_image: clean
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .


docker_image_push:
	docker push $(DOCKER_IMAGE)

docker_run:
	# are you sure docker image is up to date?
	docker run --rm $(DOCKER_IMAGE):latest version

clean:
	rm -f $(BINARY) coverage*
