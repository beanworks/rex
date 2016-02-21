.PHONY: all

all: godep

include cluster/rabbitmq.mk

godep: *.go
	GO15VENDOREXPERIMENT=0 godep go build -o rex ./

no-godep: *.go
	go get ./...
	go build -o rex ./

test:
	GO15VENDOREXPERIMENT=0 godep go vet ./...
	GO15VENDOREXPERIMENT=0 godep go test ./...

vtest:
	GO15VENDOREXPERIMENT=0 godep go vet -v ./...
	GO15VENDOREXPERIMENT=0 godep go test -v -cover ./...

clean:
	go clean ./...

cover:
	GO15VENDOREXPERIMENT=0 godep go test -coverprofile c.out ./...
	GO15VENDOREXPERIMENT=0 godep go tool cover -html=c.out
