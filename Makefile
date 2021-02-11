OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
MAKE_COLOR=\033[33;01m%-20s\033[0m

all: test

run: ## run server
	@echo "$(OK_COLOR)--> Running server$(NO_COLOR)"
	go run .

lint: ## run linters
	@echo "$(OK_COLOR)--> Running linters$(NO_COLOR)"
	golangci-lint run

test: test-unit test-acceptance ## run all tests

test-unit: ## run unit tests
	@echo "$(OK_COLOR)--> Running unit tests$(NO_COLOR)"
	go test -v --race --count=1 ./...

test-acceptance: ## run acceptance tests
	@echo "$(OK_COLOR)--> Running Acceptance tests$(NO_COLOR)"
	go test -v --tags "acceptance" --race --count=1 ./...

fmt: ## format go files
	@echo "$(OK_COLOR)--> Formatting go files$(NO_COLOR)"
	go fmt ./...

help: ## show this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(MAKE_COLOR) %s\n", $$1, $$2 } /^##@/ { printf "\n$(MAKE_COLOR)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: all run lint test test-unit test-e2e fmt clean help