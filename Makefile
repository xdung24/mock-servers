build:
	go build -o mock-servers ./...

clean:
	rm -f mock-servers

.PHONY: build clean
