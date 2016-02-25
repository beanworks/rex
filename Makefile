.PHONY: all

VENDOR_FLAG = GO15VENDOREXPERIMENT=1
GO_CMD = $(VENDOR_FLAG) godep go

all: build

include cluster/rabbitmq.mk

build: *.go
	$(GO_CMD) build -o rex ./

test:
	$(GO_CMD) vet ./...
	$(GO_CMD) test ./...

vtest:
	$(GO_CMD) vet -v ./...
	$(GO_CMD) test -v -cover ./...

clean:
	$(GO_CMD) clean ./...

cover:
	$(GO_CMD) test -coverprofile c.out ./...
	$(GO_CMD) tool cover -html=c.out
