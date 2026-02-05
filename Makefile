PROJECT=deheap
GO=go
GO_BUILD=$(GO) build
GO_TEST=$(GO) test
GO_VET=$(GO) vet
GO_BENCH=$(GO_TEST) -bench .
GIT=git
SH_TOOLS_DIR=./tools/sh
BUILD_VERSION:=$(shell cat ./VERSION 2>/dev/null || echo "0.0.0")

.PHONY: all
all: build test vet

.PHONY: build
build:
	$(GO_BUILD) ./...

.PHONY: test
test:
	$(GO_TEST) ./...

.PHONY: bench
bench:
	$(GO_BENCH) ./...

.PHONY: vet
vet:
	$(GO_VET) ./...

.PHONY: clean
clean:
	$(GO) clean -cache -testcache -fuzzcache

# Create an annotated git tag from the version in ./VERSION.
#   make tag
.PHONY: tag
tag:
	@echo "$(BUILD_VERSION)" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9_.-]+)?$$' || (echo "Error: invalid version '$(BUILD_VERSION)'"; exit 1)
	$(GIT) tag -a "$(BUILD_VERSION)" -m "Release $(BUILD_VERSION)"
	@echo "Created tag $(BUILD_VERSION)"

# Bump the major version in VERSION (resets minor and patch to 0).
#   make bump-major
.PHONY: bump-major
bump-major:
	$(SH_TOOLS_DIR)/bump-version.sh major

# Bump the minor version in VERSION (resets patch to 0).
#   make bump-minor
.PHONY: bump-minor
bump-minor:
	$(SH_TOOLS_DIR)/bump-version.sh minor

# Bump the patch version in VERSION.
#   make bump-patch
.PHONY: bump-patch
bump-patch:
	$(SH_TOOLS_DIR)/bump-version.sh patch
