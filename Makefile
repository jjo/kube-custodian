NAME ?= kube-custodian
LDFLAGS ?= -ldflags="-s -w -X main.version=$(VERSION) -X main.revision=$(REVISION)"
VERSION ?= master
GOARCH ?= amd64

GOSRC = $(shell find *.go pkg cmd -name '*.go')

ifeq ($(GOARCH), arm)
DOCKERFILE_SED_EXPR?=s,FROM alpine:,FROM multiarch/alpine:armhf-v,
DOCKER_IMG_FULL=$(DOCKER_IMG):arm-$(VERSION)
else ifeq ($(GOARCH), arm64)
DOCKERFILE_SED_EXPR?=s,FROM alpine:,FROM multiarch/alpine:aarch64-v,
DOCKER_IMG_FULL=$(DOCKER_IMG):arm64-$(VERSION)
else
DOCKERFILE_SED_EXPR?=
DOCKER_IMG_FULL=$(DOCKER_IMG):$(VERSION)
endif
GOPKGS = $(shell glide novendor)

DOCKER_REPO ?= quay.io
DOCKER_IMG ?= $(DOCKER_REPO)/jjo/kube-custodian


all: build

build: bin/$(NAME)

bin/$(NAME): $(GOSRC)
	GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(NAME)

lint:
	golint $(GOPKGS)

test:
	go test -v -cover $(GOPKGS)

clean:
	rm -fv bin/$(NAME)


docker-build: Dockerfile.$(GOARCH).run
	docker build --build-arg SRC_TAG=$(VERSION) --build-arg ARCH=$(GOARCH) -t $(DOCKER_IMG_FULL) -f $(^) .

Dockerfile.%.run: Dockerfile
	@sed -e "$(DOCKERFILE_SED_EXPR)" Dockerfile > $(@)


docker-push:
	docker push $(DOCKER_IMG)

docker-clean:
	docker image rm $(DOCKER_IMG)


multiarch-setup:
	docker run --rm --privileged multiarch/qemu-user-static:register
	dpkg -l qemu-user-static

.PHONY: all build lint test clean docker-build docker-push docker-clean
