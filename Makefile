.PHONY: build test install snapshot release-check

build:
	go build -o shipkit ./cmd/shipkit

test:
	go test ./...

install:
	go install ./cmd/shipkit

snapshot:
	goreleaser release --snapshot --clean

release-check:
	goreleaser check
