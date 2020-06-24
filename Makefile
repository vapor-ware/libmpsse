#
# libmpsse
#

# The make install below requires sudo.
.PHONY: install
install: ## Install the libmpsse package with python disabled
	cd src ; ./configure --disable-python
	cd src ; make
	cd src ; make install
	cd src ; make distclean

.PHONY: clean
clean:  ## Remove temporary files
	go clean -v || exit

.PHONY: build
build: install ## Install the libmpsse package and run 'go build'
	go build

.PHONY: fmt
fmt:  ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file" || exit; done

.PHONY: lint
lint: ## Lint the Go source code
	golint .

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
