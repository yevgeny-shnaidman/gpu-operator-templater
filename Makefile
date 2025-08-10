
.PHONY: templater
templater: operator-sdk ## Build manager binary.
	go build -ldflags="-X main.Version=$(PROJECT_VERSION) -X main.GitCommit=$(GIT_COMMIT)" -o $@ ./cmd

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

OPERATOR_SDK = $(shell pwd)/internal/operator_sdk/binaries/operator-sdk
.PHONY: operator-sdk
operator-sdk:
	@if [ ! -f ${OPERATOR_SDK} ]; then \
                set -e ;\
                echo "Downloading ${OPERATOR_SDK}"; \
                mkdir -p $(dir ${OPERATOR_SDK}) ;\
                curl -Lo ${OPERATOR_SDK} 'https://github.com/operator-framework/operator-sdk/releases/download/v1.40.0/operator-sdk_linux_amd64'; \
                chmod +x ${OPERATOR_SDK}; \
        fi
