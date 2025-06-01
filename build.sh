#!/bin/bash

set -ex

# Get version information
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X 'rayyanriaz/tool-version-manager/cmd/tools-manager.Version=$VERSION' \
         -X 'rayyanriaz/tool-version-manager/cmd/tools-manager.Commit=$COMMIT' \
         -X 'rayyanriaz/tool-version-manager/cmd/tools-manager.BuildDate=$BUILD_DATE'"

# Build for current platform
echo "Building for $(go env GOOS)/$(go env GOARCH)..."
mkdir -p bin
go build -ldflags "$LDFLAGS" -o bin/tvm .

echo "Build complete: bin/tvm"
