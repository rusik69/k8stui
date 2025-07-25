# K8sTUI - Kubernetes Terminal UI

A lightweight, keyboard-driven terminal UI for managing Kubernetes clusters, inspired by k9s.

## Features

- View and navigate Kubernetes namespaces
- View pods within selected namespaces
- View containers within selected pods
- Hotkey-based navigation
- Delete resources with confirmation
- Real-time log viewing

## Prerequisites

- Go 1.22
- Access to a Kubernetes cluster (via kubeconfig)

## Installation

```bash
# Clone the repository
git clone https://github.com/rusik69/k8stui.git
cd k8stui

# Build the application
make build

# Install to your GOPATH
make install
```

## Usage

```bash
# Run the application
k8stui
```

### Hotkeys

- `TAB`/`Shift+TAB`: Navigate between panels
- `ENTER`: Select item
- `Ctrl+D`: Delete selected resource (with confirmation)
- `Q`: Quit application
- `↑/↓/←/→`: Scroll through content

## Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Run the application in development mode
make run

# Build for different platforms
make build-all  # Builds for all supported platforms
make build-darwin-arm64  # Build for Apple Silicon
make build-linux-amd64   # Build for Linux
```

## GitHub Actions

This project includes GitHub Actions for continuous integration:
- **CI**: Runs tests and builds on Go 1.22
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
