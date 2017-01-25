Experimental project using go-kit
------

#### Build
``GOOS=linux make build``
``make container``

#### Run
``docker run --rm -p 9001:9001 --name posts mboard/posts:latest
``docker run --rm -p 8080:8080 --link posts:posts mboard/api:latest``

#### Deploy (from macOS) using Kubernetes
``find -E . -regex  "(.*)+\-(controller|service|deployment).(json|yml)" | xargs -L 1 kubectl create -f``
> Deploy on linux: omit the -E flag