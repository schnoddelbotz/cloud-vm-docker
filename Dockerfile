FROM golang:1.14.2 AS builder

WORKDIR /src/github.com/schnoddelbotz/cloud-vm-docker
COPY . .

RUN make test clean build CGO_ENABLED=0

FROM alpine
RUN apk add ca-certificates
COPY --from=builder /src/github.com/schnoddelbotz/cloud-vm-docker/cloud-vm-docker /bin/cloud-vm-docker

ENTRYPOINT ["/bin/cloud-vm-docker"]
CMD []
