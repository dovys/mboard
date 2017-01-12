SHELL = /bin/bash
BINARY = mboard
TEST_DIRECTORIES = ./handlers ./services

install:
	@which glide > /dev/null || { echo "glide not found. run brew install glide"; exit 1; }
	@glide install

testv:
	@go test -v $(TEST_DIRECTORIES)

test:
	@go test $(TEST_DIRECTORIES)

build:
	@go build -o $(BINARY)
	@echo "Binary:" $(CURDIR)/$(BINARY)

clean:
	@if [ -f $(BINARY) ]; then rm $(BINARY); fi

.PHONY: clean install