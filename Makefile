#
# libmpsse
#

IMAGE   := vaporio/libmpsse-base
VERSION := 2.0

.PHONY: docker
docker: ## Build the docker images
	# building the base image
	docker build -f dockerfile/base.Dockerfile \
		-t $(IMAGE):latest \
		-t $(IMAGE):$(VERSION) .
	# building the builder image
	docker build -f dockerfile/build.Dockerfile \
		-t $(IMAGE):builder .

.PHONY: install
install: ## Install the libmpsse package with python disabled
	cd src ; make distclean
	cd src ; ./configure --disable-python
	cd src ; make
	cd src ; make install

.PHONY: build
build: install ## Install the libmpsse package and run 'go build'
	go build

.PHONY: lint
lint: ## Lint the Go source code
	golint .

.PHONY: push-dockerhub
push-dockerhub: ## Push the images to DockerHub
	docker push $(IMAGE):lastest
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):builder

.PHONY: version
version: ## Print the version
	@echo "$(VERSION)"

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
