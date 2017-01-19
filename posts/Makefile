SHELL = /bin/bash
MAKEFLAGS+=-s
BINARY = build/service
TEST_DIRECTORIES = ./ 
.DEFAULT_GOAL := build

build:
	go build -o $(BINARY) cmd/main.go
	echo $(BINARY)

install:
	which glide > /dev/null || { echo "glide not found. run brew install glide"; exit 1; }
	glide install

testv:
	go test -v $(TEST_DIRECTORIES)

test:
	go test $(TEST_DIRECTORIES)

clean:
	if [ -f $(BINARY) ]; then rm $(BINARY); fi

container:
	which docker > /dev/null || { echo "Docker not found."; exit 1; }
	docker build -t mboard/posts:latest .
	docker tag mboard/posts:latest gcr.io/mboard-kub/posts

proto:
	which protoc > /dev/null || { echo "protoc not found."; exit 1; }
	protoc pb/*.proto --go_out=plugins=grpc:.

.PHONY: clean install test testv proto container build
