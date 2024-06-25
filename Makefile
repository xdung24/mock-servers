GO_OS := $(shell go env GOOS)
GO_ARCH := $(shell go env GOARCH)

build:
	CGO_ENABLED=0 GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -o mock-servers -ldflags "-X main.Goos=$(GO_OS) -X main.Goarch=$(GO_ARCH)" 

clean:
	rm -f mock-servers

.PHONY: build clean