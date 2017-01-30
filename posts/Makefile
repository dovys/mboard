SHELL = /bin/bash
MAKEFLAGS+=-s
BINARY = build/service
TEST_DIRECTORIES = ./ 
.DEFAULT_GOAL := build

build:
	go build -o $(BINARY) cmd/main.go
	echo $(BINARY)

testv:
	go test -v $(TEST_DIRECTORIES)

test:
	go test $(TEST_DIRECTORIES)

clean:
	if [ -f $(BINARY) ]; then rm $(BINARY); fi

container:
	which docker > /dev/null || { echo "Docker not found."; exit 1; }
	docker build -t mboard/posts:latest -t gcr.io/mboard-kub/posts .

proto:
	which protoc > /dev/null || { echo "protoc not found."; exit 1; }
	protoc pb/*.proto --go_out=plugins=grpc:.

.PHONY: clean test testv proto container build
