REPO=ecr-token-refresh

.PHONY: build docker

build:
	docker run --rm -v $$(pwd):/go/src/github.com/skuid/$(REPO) \
		-w /go/src/github.com/skuid/$(REPO) \
		golang:1.8 go build -v -a -tags netgo -installsuffix netgo -ldflags '-w'

docker: build
	docker build -t quay.io/skuid/$(REPO) .

push:
	docker push quay.io/skuid/$(REPO)