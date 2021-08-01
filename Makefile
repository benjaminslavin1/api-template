.PHONY: test # Runs unit tests
test:
	@go test ./...

.PHONY: lint # Report all linting issues for repo
lint:
	@golangci-lint run