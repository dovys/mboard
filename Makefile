SHELL = /bin/bash
BINARY = mboard

install:
	@which glide > /dev/null || { echo "glide not found. run brew install glide"; exit 1; }
	@glide install

build:
	@go build -o $(BINARY)
	@echo "Binary:" $(CURDIR)/$(BINARY)

clean:
	@if [ -f $(BINARY) ]; then rm $(BINARY); fi

.PHONY: clean install