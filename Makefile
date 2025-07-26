.PHONY: fmt lint

fmt:
	@golines -w .
	@gofumpt -w .

lint:
	@golangci-lint run --fix