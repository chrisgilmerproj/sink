NO_V_VERSION=0.0.1
VERSION=v$(NO_V_VERSION)
DEV_VERSION=$(VERSION)-devel

LDFLAGS=-X main.version=$(VERSION) -extldflags=-static -s -w

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[\/a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#
# Commands
#

.PHONY: build
build: ## Build the binary
	CGO_ENABLED=0 go build -installsuffix cgo -ldflags "-X main.version=$(DEV_VERSION) -extldflags=-static -s -w" -o bin/sink main.go

.PHONY: install
install: build ## Install the binary
	sudo cp bin/sink /usr/local/bin/sink

.PHONY: run
run: ## Run the code
	go run main.go

.PHONY: test
test: ## Run the tests
	go test ./...

.PHONY: test_verbose
test_verbose: ## Run the tests in verbose mode with coverage
	go test ./... -v -cover -count=1

.PHONY: clean
clean: ## Clean the build files
	rm -rf bin/*
