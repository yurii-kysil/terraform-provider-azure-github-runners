.PHONY: build test clean install release

# Build the provider
build:
	go build -o terraform-provider-azure-github-runners

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f terraform-provider-azure-github-runners

# Create a release (requires goreleaser)
release:
	goreleaser release --clean

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate documentation
docs:
	cd tools && go generate ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
