GO ?= go

.PHONY: test lint check run

test:
	$(GO) test ./...

lint:
	$(GO) fmt ./...
	$(GO) vet ./...

check: lint test

run:
	$(GO) run ./cmd/gaussgo
