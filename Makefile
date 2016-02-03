.PHONY: all

all: godep

include cluster/rabbitmq.mk

godep: *.go
	godep go build -o rex ./

no-godep: *.go
	go get ./...
	go build -o rex ./

test:
	godep go vet ./...
	godep go test ./...

vtest:
	godep go vet -v ./...
	godep go test -v -cover ./...

clean:
	godep go clean ./...

cover:
	godep go test -coverprofile c.out ./...
	godep go tool cover -html=c.out
