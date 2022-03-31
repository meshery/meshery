IMAGE?=meshery/meshery-docker-desktop-extension

BUILDER=buildx-multi-arch

STATIC_FLAGS=CGO_ENABLED=0
LDFLAGS="-s -w"
GO_BUILD=$(STATIC_FLAGS) go build -trimpath -ldflags=$(LDFLAGS)

GIT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_STRIPPED_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1` | cut -c 2-)
INFO_COLOR = \033[0;36m
NO_COLOR   = \033[m

bin: ## Build the binary for the current platform
	@echo "$(INFO_COLOR)Building...$(NO_COLOR)"
	$(GO_BUILD) -o bin/service ./vm

extension: ## Build service image to be deployed as a desktop extension
	docker build --tag=$(IMAGE) --build-arg GIT_VERSION=$(GIT_VERSION) --build-arg GIT_STRIPPED_VERSION=$(GIT_STRIPPED_VERSION) .

prepare-buildx: ## Create buildx builder for multi-arch build, if not exists
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host

push-extension: prepare-buildx ## Build & Upload extension image to hub. Do not push if tag already exists: make push-extension tag=0.1
	docker pull $(IMAGE):$(tag) && echo "Failure: Tag already exists" || docker buildx build --push --builder=$(BUILDER) --platform=linux/amd64,linux/arm64 --build-arg TAG=${tag)} --tag=$(IMAGE):$(tag) .

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'

.PHONY: bin extension push-extension help
