# Build Documentation

This document describes how to build, test, and produce release artifacts for the `go-api-1775281337` gateway service.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/) installed and on your `PATH`
- `git` for cloning the repository
- `sha256sum` (Linux) or `shasum -a 256` (macOS) for checksum verification

## Repository Structure

```
.
├── gateway/          # Main application (HTTP gateway service)
│   ├── main.go       # Entry point
│   └── handler/      # HTTP handlers
├── go.mod            # Go module definition
├── go.sum            # Dependency lock file
├── .github/
│   └── workflows/
│       └── go-ci.yml # CI/CD pipeline definition
└── BUILD.md          # This file
```

## Building Locally

### Quick Build

```bash
# Build the gateway binary to the current directory
go build -o gateway ./gateway/...
```

### Production Build

The production build strips debug symbols (`-s -w`) to reduce binary size and embeds version metadata:

```bash
# Create output directory
mkdir -p dist

# Build with version metadata
go build \
  -ldflags="-s -w -X main.version=$(git rev-parse --short HEAD) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o dist/gateway \
  ./gateway/...
```

### Cross-Compilation

To build for a different OS/architecture:

```bash
# Linux AMD64 (e.g., for deployment to a Linux server from macOS)
GOOS=linux GOARCH=amd64 go build -o dist/gateway-linux-amd64 ./gateway/...

# Linux ARM64 (e.g., for AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o dist/gateway-linux-arm64 ./gateway/...
```

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Tests with Race Detector

Always run with `-race` before submitting a PR to catch data races:

```bash
go test -race ./...
```

### Run Tests with Coverage

```bash
# Generate coverage profile
go test -race -coverprofile=coverage.out ./...

# View coverage summary in terminal
go tool cover -func=coverage.out

# Open interactive HTML coverage report in browser
go tool cover -html=coverage.out
```

### Run a Specific Package or Test

```bash
# Test a specific package
go test ./gateway/handler/...

# Run a specific test by name
go test -run TestHelloHandler ./gateway/handler/...
```

## Code Quality

```bash
# Run the Go static analyser
go vet ./...

# Format code (modifies files in place)
gofmt -w .

# Check formatting without modifying files
gofmt -l .
```

## Generating Checksums

After building, generate a SHA-256 checksum to verify binary integrity:

```bash
# Linux
cd dist && sha256sum gateway > gateway.sha256

# macOS
cd dist && shasum -a 256 gateway > gateway.sha256

# Verify the checksum
sha256sum --check gateway.sha256   # Linux
shasum -a 256 --check gateway.sha256  # macOS
```

## CI/CD Artifacts

### What CI Produces

The GitHub Actions workflow (`.github/workflows/go-ci.yml`) produces the following artifacts on every successful build:

| Artifact Name | Contents | Retention |
|---|---|---|
| `gateway-binary-<sha>` | Compiled binary (`dist/gateway`) + checksum (`dist/gateway.sha256`) | 30 days |
| `gateway-checksum-<sha>` | SHA-256 checksum file only (`dist/gateway.sha256`) | 30 days |
| `coverage-report` | Go test coverage profile (`coverage.out`) | 7 days |

### Downloading Artifacts

Artifacts can be downloaded from the **Actions** tab in GitHub:

1. Navigate to the repository on GitHub
2. Click **Actions** in the top navigation
3. Select the workflow run of interest
4. Scroll to the **Artifacts** section at the bottom of the run summary
5. Click the artifact name to download a ZIP archive

Alternatively, use the [GitHub CLI](https://cli.github.com/):

```bash
# List artifacts for a workflow run
gh run list --workflow=go-ci.yml

# Download artifacts from a specific run
gh run download <run-id>
```

### Verifying a Downloaded Artifact

```bash
# After downloading and extracting the artifact ZIP:
sha256sum --check gateway.sha256   # Linux
shasum -a 256 --check gateway.sha256  # macOS
```

### Artifact Naming Convention

Artifacts are named with the full Git commit SHA suffix (e.g., `gateway-binary-abc1234...`) to ensure:
- Each build's artifacts are uniquely identifiable
- Artifacts can be traced back to a specific commit
- No ambiguity when multiple workflow runs exist for the same branch

## Running the Service Locally

```bash
# Build and run in one step
go run ./gateway/...

# Or run the compiled binary
./dist/gateway

# Configure the port (default: 8080)
PORT=9090 ./dist/gateway
```

The service exposes the following endpoints:

| Method | Path | Description |
|---|---|---|
| `GET` | `/hello` | Hello handler |

## Dependency Management

```bash
# Download all dependencies
go mod download

# Verify integrity of downloaded modules
go mod verify

# Tidy (remove unused, add missing)
go mod tidy
```
