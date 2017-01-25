SHELL = /bin/bash
MAKEFLAGS+=-s
.DEFAULT_GOAL := install

install:
	which dep > /dev/null || { echo "github.com/golang/dep not found"; exit 1; }
	dep ensure