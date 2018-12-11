#
# libmpsse
#

.PHONY: install
install: ## Install the libmpsse package with python disabled
	cd src ; ./configure --disable-python
	cd src ; make
	cd src ; make install
	cd src ; make distclean

.PHONY: build
build: install ## Install the libmpsse package and run 'go build'
	go build

.PHONY: lint
lint: ## Lint the Go source code
	golint .

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
