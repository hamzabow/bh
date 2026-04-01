# ============================================================================
# bh - Base Converter TUI Makefile
# ============================================================================
# Self-documenting Makefile with auto-generated help from inline comments.
#
# Quick Start:
#   make help    # Show all available commands
#   make build   # Build the binary
#   make run     # Build and run the TUI
# ============================================================================

BINARY := bh

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show all available targets with descriptions
	@echo ""
	@echo "  \033[1mbh\033[0m — Terminal UI for converting between number bases"
	@echo "  (decimal, hexadecimal, octal, binary)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""

##@ Build

.PHONY: build
build: ## Build the binary
	go build -o $(BINARY)

.PHONY: run
run: build ## Build and run the TUI
	./$(BINARY)

.PHONY: install
install: ## Install the binary to GOBIN
	go install .

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BINARY)

# ##@ Quality (uncomment when ready)
#
# .PHONY: lint
# lint: ## Run golangci-lint (requires: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
# 	golangci-lint run
#
# .PHONY: test
# test: ## Run tests
# 	go test ./...
