SHELL = /bin/bash
MAKEFLAGS+=-s
BINARY = build/api
TEST_DIRECTORIES = ./handlers ./services
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
	docker build -t mboard/api:latest -t gcr.io/mboard-kub/api .

.PHONY: clean test testv build container