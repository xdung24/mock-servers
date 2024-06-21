GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

build:
	go build -o mock-servers -ldflags "-X main.Goos=$(GOOS) -X main.Goarch=$(GOARCH)" ./...

clean:
	rm -f mock-servers

.PHONY: build clean