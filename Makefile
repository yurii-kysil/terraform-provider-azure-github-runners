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

# Install the provider locally
install: build
	mkdir -p ~/.terraform.d/plugins/local/yurii-kysil/azure-github-runners/1.0.0/linux_amd64
	cp terraform-provider-azure-github-runners ~/.terraform.d/plugins/local/yurii-kysil/azure-github-runners/1.0.0/linux_amd64/

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
