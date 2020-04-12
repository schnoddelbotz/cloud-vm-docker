
BINARY := cloud-task-zip-zap

VERSION := $(shell git describe --tags | cut -dv -f2)
DOCKER_IMAGE := schnoddelbotz/cloud-task-zip-zap
LDFLAGS := -X github.com/schnoddelbotz/cloud-task-zip-zap/cmd.AppVersion=$(VERSION) -w

GO_SOURCES := */*.go */*/*.go */*/*/*.go


build: $(BINARY)
	# running: $(BINARY) version
	@./$(BINARY) version

$(BINARY): $(GO_SOURCES)
	# building cloud-task-zip-zap
	go build -v -o $(BINARY) -ldflags='-w -s $(LDFLAGS)' ./ctzz/main.go

all_local: clean test build

all_docker: clean test docker_image_prod

release: all_docker docker_image_push


test:
	# golint -set_exit_status ./...
	golint ./...
	go vet ./...
	go test -ldflags='-w -s $(LDFLAGS)' ./...

coverage: clean
	PROVIDER=MEMORY go test -coverprofile=coverage.out -coverpkg=./... -ldflags='-w -s $(LDFLAGS)' ./...
	go tool cover -html=coverage.out

deploy_gcp:
	make -j2 deploy_cfn_http deploy_cfn_pubsub

deploy_cfn_http:
	gcloud functions deploy CloudTaskZipZap --region=europe-west1 --runtime go113 \
 		--trigger-http --allow-unauthenticated \
 		--set-env-vars=CTZZ_DATASTORE_COLLECTION=cloud-task-zip-zap-test

deploy_cfn_pubsub:
	gcloud functions deploy CloudTaskZipZapProcessor --region=europe-west1 --runtime go113 \
     		--trigger-topic=ctzz-task-queue --allow-unauthenticated \
     		--set-env-vars=CTZZ_TOPIC=ctzz-task-queue

docker_image:
	docker build -f Docker/cloud-task-zip-zap.Dockerfile -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .


docker_image_push:
	docker push $(DOCKER_IMAGE)

docker_run:
	# are you sure docker image is up to date?
	docker run --rm $(DOCKER_IMAGE):latest version

clean:
	rm -f $(BINARY) coverage*
