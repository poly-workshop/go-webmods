.PHONY: fmt lint

fmt:
	@gofumpt -w .

lint:
	@golangci-lint run --fix