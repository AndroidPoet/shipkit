.PHONY: build test cover cover-html lint fmt fmt-check install snapshot release-check tidy

build:
	go build -o shipkit ./cmd/shipkit

test:
	go test ./...

# Run the full suite with an atomic coverage profile.
cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

# Produce a browsable HTML coverage report.
cover-html: cover
	go tool cover -html=coverage.out -o coverage.html
	@echo "Wrote coverage.html"

# Static analysis. Requires golangci-lint (https://golangci-lint.run).
lint:
	golangci-lint run ./...

# Format the codebase with gofumpt (stricter gofmt).
fmt:
	gofumpt -w .

# Fail if any file is not gofumpt-clean.
fmt-check:
	@unformatted="$$(gofumpt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files are not gofumpt-clean:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

tidy:
	go mod tidy

install:
	go install ./cmd/shipkit

snapshot:
	goreleaser release --snapshot --clean

release-check:
	goreleaser check
