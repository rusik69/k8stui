# Kubernetes TUI (kUI)

A simple Kubernetes text-based user interface written in Golang, inspired by k9s.

## Features

- View and navigate Kubernetes namespaces
- View pods within selected namespaces
- View containers within selected pods
- Hotkey-based navigation:

## Build Instructions

### Prerequisites
- Go 1.19 or higher
- Access to a Kubernetes cluster (via kubeconfig)

### Building

```bash
# Build the application
make build

# Run tests
make test

# Build for Linux (cross-compilation)
make build-linux

# Clean build artifacts
make clean

# Install dependencies
make deps
```

### Development

```bash
# Run the application in development
make run

# Run tests with verbose output
go test -v ./...
```

## GitHub Actions

This project includes GitHub Actions for continuous integration:
- **CI**: Runs tests and builds on multiple Go versions (1.19, 1.20, 1.21)
- **Build**: Creates binaries for multiple platforms
- **Artifacts**: Uploads built binaries for releases

## Features
  - `j` or `↓`: Move down
  - `k` or `↑`: Move up
  - `h` or `←`: Go back
  - `r`: Refresh current view
  - `q` or `Ctrl+C`: Quit
  - `Enter`: Select item

## Installation

```bash
go install github.com/yourusername/kui@latest
```

## Usage

```bash
kui
```

The application will automatically connect to your Kubernetes cluster using:
1. In-cluster configuration (if running inside a pod)
2. Local kubeconfig file (if running locally)

## Requirements

- Go 1.21 or higher
- Access to a Kubernetes cluster
- Proper kubectl/kubeconfig configuration

## Key Bindings

| Key | Action |
|-----|--------|
| `j`/`↓` | Move down |
| `k`/`↑` | Move up |
| `h`/`←` | Go back |
| `r` | Refresh |
| `q`/`Ctrl+C` | Quit |
| `Enter` | Select item |

## Development

To run the application locally:

```bash
go run main.go
```
