# ==== Test targets ====

##@ Test

.PHONY: test
test: test-unit ## Run all tests.

.PHONY: test-unit
test-unit: ## Run all unit tests.
	@GOEXPERIMENT=jsonv2 go test -count=1 -v ./...
