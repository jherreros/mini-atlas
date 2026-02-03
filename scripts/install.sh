#!/usr/bin/env bash
set -euo pipefail

REPO="jherreros/shoulders"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

if [[ "$OS" != "darwin" && "$OS" != "linux" ]]; then
  echo "Unsupported OS: $OS"
  exit 1
fi

VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4)
if [ -z "$VERSION" ]; then
  echo "Failed to fetch latest version"
  exit 1
fi

# Remove 'v' prefix for the binary filename
BINARY_VERSION=${VERSION#v}
TARBALL="shoulders_${BINARY_VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARBALL}"

TMP_DIR=$(mktemp -d)
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

echo "Downloading $URL..."
if ! curl -fsSL "$URL" | tar -xz -C "$TMP_DIR"; then
  echo "Failed to download or extract archive"
  exit 1
fi

if [[ ! -f "$TMP_DIR/shoulders" ]]; then
  echo "Binary not found in archive"
  exit 1
fi

if [[ ! -w "$INSTALL_DIR" ]]; then
  echo "Installing to $INSTALL_DIR requires sudo"
  sudo mv "$TMP_DIR/shoulders" "$INSTALL_DIR/shoulders"
else
  mv "$TMP_DIR/shoulders" "$INSTALL_DIR/shoulders"
fi

chmod +x "$INSTALL_DIR/shoulders"

echo "Installed shoulders to $INSTALL_DIR/shoulders"
