BINARY := terraform-provider-oack
VERSION := 0.1.0
OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR := ~/.terraform.d/plugins/registry.terraform.io/oack-io/oack/$(VERSION)/$(OS_ARCH)

.PHONY: default build install clean lint test testacc fmt vet

default: build

build:
	go build -o $(BINARY)

install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY) $(PLUGIN_DIR)/$(BINARY)_v$(VERSION)

clean:
	rm -f $(BINARY)

# ── Code quality ──────────────────────────────────────────────────────────────

fmt:
	gofmt -w -s .

vet:
	go vet ./...

lint:
	golangci-lint run --tests=false ./...

# ── Tests ─────────────────────────────────────────────────────────────────────

test:
	go test ./... -count=1

testacc:
	TF_ACC=1 go test ./internal/acctest/ -v -timeout 30m $(TESTARGS)
