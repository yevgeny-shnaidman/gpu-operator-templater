
.PHONY: templater
templater: $(shell find -name "*.go") go.mod ## go.sum  ## Build manager binary.
	go build -ldflags="-X main.Version=$(PROJECT_VERSION) -X main.GitCommit=$(GIT_COMMIT)" -o $@ ./cmd

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...
