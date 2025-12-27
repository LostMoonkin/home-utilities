#!/bin/bash

# Configuration
BINARY_NAME="homeserver"
BUILD_DIR="build"

# Automatically detect current platform
TARGET_OS=$(go env GOOS)
TARGET_ARCH=$(go env GOARCH)

echo "Detected Platform: $TARGET_OS/$TARGET_ARCH"

# 1. Clean build directory
if [ -d "$BUILD_DIR" ]; then
    echo "Cleaning existing $BUILD_DIR directory..."
    rm -rf "$BUILD_DIR"
fi

# 2. Create fresh build directory
echo "Creating fresh directory: $BUILD_DIR"
mkdir -p "$BUILD_DIR"

# 3. Execute compilation
# CGO_ENABLED=0 ensures a static binary
# -ldflags="-s -w" reduces binary size
echo "Starting build process..."

CGO_ENABLED=0 \
GOOS=$TARGET_OS \
GOARCH=$TARGET_ARCH \
go build -ldflags="-s -w" -trimpath -o "$BUILD_DIR/$BINARY_NAME" .

if [ $? -eq 0 ]; then
    echo "Build successful: $BUILD_DIR/$BINARY_NAME"
else
    echo "Error: Compilation failed"
    exit 1
fi

# 4. Copy .env file to the build directory
if [ -f ".env" ]; then
    cp .env "$BUILD_DIR/"
    echo ".env file copied to $BUILD_DIR/"
else
    echo "Warning: .env file not found in current directory, skipping copy."
fi

echo "Process finished successfully."
ls -lh "$BUILD_DIR"
