IMAGE_REPO ?= localhost:5000/imageswap
IMAGE_NAME ?= imageswap

GIT_HOST =? github.com

PWD := $(shell pwd)
BASE_DIR := $(shell basename $(PWD))

# Keep an existing GOPATH, make a private one if it is undefined
GOPATH_DEFAULT := $(PWD)/.go
export GOPATH ?= $(GOPATH_DEFAULT)
TESTARGS_DEFAULT := "-v"
export TESTARGS ?= $(TESTARGS_DEFAULT)
DEST := $(GOPATH)/src/$(GIT_HOST)/$(BASE_DIR)
IMAGE_TAG ?= $(shell date +v%Y%m%d)-$(shell git describe --match=$(git rev-parse --short=8 HEAD) --tags --always --dirty)

LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
    TARGET_OS ?= linux
    XARGS_FLAGS="-r"
else ifeq ($(LOCAL_OS),Darwin)
    TARGET_OS ?= darwin
    XARGS_FLAGS=
else
    $(error "This system's OS $(LOCAL_OS) isn't recognized/supported")
endif

all: fmt lint test build image

ifeq (,$(wildcard go.mod))
ifneq ("$(realpath $(DEST))", "$(realpath $(PWD))")
    $(error Please run 'make' from $(DEST). Current directory is $(PWD))
endif
endif

############################################################
# format section
############################################################

fmt:
	@echo "Run go fmt..."
	@go fmt $(PWD)

############################################################
# lint section
############################################################

lint:
	@echo "Runing the golangci-lint..."
	@golangci-lint run ./...

############################################################
# test section
############################################################

test:
	@echo "Running the tests for $(IMAGE_NAME)..."
	@go test $(TESTARGS) ./...

test-deployment: image
	@echo "Deploying webhook and webhook configuration on local kind cluster"
	@./internal/deploy/setup.sh
	@echo ""
	@echo ""
	@echo "Test image hostname mutating by running 'kubectl run alpine --image=alpine --restart=Never -n default --command -- sleep infinity' and confirming your ECR hostname appears as the 'image:' hostname"

############################################################
# build section
############################################################

build:
	@echo "Building the $(IMAGE_NAME) binary..."
	@CGO_ENABLED=0 go build -o build/_output/bin/$(IMAGE_NAME) main.go

build-linux:
	@echo "Building the $(IMAGE_NAME) binary for linux..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/_output/linux/bin/$(IMAGE_NAME) main.go

############################################################
# image section
############################################################

image: build-image push-image

build-image: build-linux
	@echo "Building the docker image: $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker build -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) .

push-image: build-image
	@echo "Pushing the docker image for $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)"
	@docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

############################################################
# clean section
############################################################
clean:
	@rm -rf build/_output

.PHONY: all fmt lint check test build image clean
