[![Build Status](https://travis-ci.org/dovys/mboard.svg?branch=master)](https://travis-ci.org/dovys/mboard)

Experimental project using go-kit
------

#### Build
```sh
GOOS=linux make build
make container
```

#### Run
```sh
docker run --rm -p 9001:9001 --name posts mboard/posts:latest
docker run --rm -p 8080:8080 --link posts:posts mboard/api:latest
```

#### Deploy (from macOS) using Kubernetes
```sh
find -E . -regex  "(.*)+\-(controller|service|deployment).(json|yml)" | xargs -L 1 kubectl create -f
```
> Deploy from linux: omit the -E flag

#### Todo
- Proper logging, tracing
- Metrics on Prometheus
- Auth
- TLS
- Better coverage
- ~~CI builds~~
- To commit vendor or not to commit vendor
