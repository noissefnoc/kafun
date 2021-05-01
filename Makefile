.DEFAULT_GOAL := help

deps: ## Install dependencies
	go mod download

test: deps ## Run tests
	go test ./...

build: deps ## build binary
	mkdir bin && go build -o bin/kafun cmd/kafun/main.go

clean:
	rm -rf bin

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'

.PHONY: deps test build clean help
