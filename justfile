# Tasks for pea (Go CLI)
set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

bin := "pea"

# Default lists available tasks
default: help

help:
  @just --list

# Build a local binary
build:
  go build -o {{bin}} .

# Run the app with optional args: `just run -- --flag value`
run *args:
  go run . {{args}}

# Format source files
fmt:
  go fmt ./...

# Vet code for suspicious constructs
vet:
  go vet ./...

# Lint: vet + staticcheck if available
lint:
  @echo "Running vet"
  go vet ./...
  @echo "Running staticcheck (if available)"
  if command -v staticcheck >/dev/null 2>&1; then staticcheck ./...; else echo "staticcheck not installed, skipping"; fi

# Tidy module dependencies
tidy:
  go mod tidy

# Run tests (add e2e tests under ./e2e when available)
test:
  go test ./...

# End-to-end tests, if ./e2e exists
e2e:
  if [ -d "e2e" ]; then go test -v ./e2e; else echo "No e2e tests yet"; fi

# Aggregate checks required by the process
check: fmt vet test
  @echo "check: fmt, vet, test passed"

# Pre-commit hook runner
pre-commit: check
  @echo "pre-commit checks complete"

# Clean build artifacts
clean:
  rm -f {{bin}}