.PHONY: all

GO_CMD = godep go
PKG = $$(go list ./... | grep -v /vendor/)

all: build

include cluster/rabbitmq.mk

build: *.go
	$(GO_CMD) build -race -o rex ./

test:
	$(GO_CMD) vet $(PKG)
	$(GO_CMD) test -race $(PKG)

vtest:
	$(GO_CMD) vet -v $(PKG)
	$(GO_CMD) test -race -v -cover $(PKG)

clean:
	$(GO_CMD) clean $(PKG)

cover:
	@echo "mode: count" > c.out
	@for pkg in $(PKG); do \
		$(GO_CMD) test -coverprofile c.out.tmp $$pkg; \
		tail -n +2 c.out.tmp >> c.out; \
	done
	$(GO_CMD) tool cover -html=c.out

rel: release

release:
	./dist/release.sh
