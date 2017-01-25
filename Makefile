SHELL = /bin/bash
MAKEFLAGS+=-s
.DEFAULT_GOAL := install

install:
	which dep > /dev/null || { echo "github.com/golang/dep not found"; exit 1; }
	dep ensure

build:
	ls -d */ | grep -v vendor | xargs -L 1 make clean build -C

container:
	ls -d */ | grep -v vendor | xargs -L 1 make container -C
