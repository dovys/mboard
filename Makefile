SHELL = /bin/bash
MAKEFLAGS+=-s
BINARY = mboard
TEST_DIRECTORIES = ./handlers ./services

$(BINARY):
	go build -o $(BINARY)
	echo "Binary:" $(CURDIR)/$(BINARY)

install:
	which glide > /dev/null || { echo "glide not found. run brew install glide"; exit 1; }
	glide install

testv:
	go test -v $(TEST_DIRECTORIES)

test:
	go test $(TEST_DIRECTORIES)

clean:
	if [ -f $(BINARY) ]; then rm $(BINARY); fi

.PHONY: clean install test testv