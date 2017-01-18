SHELL = /bin/bash
MAKEFLAGS+=-s
BINARY = service
TEST_DIRECTORIES = ./ 

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

proto:
	which protoc > /dev/null || { echo "protoc not found."; exit 1; }
	protoc pb/*.proto --go_out=plugins=grpc:.

.PHONY: clean install test testv proto
