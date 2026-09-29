#!/bin/sh
set -e

REPO="holasoymalva/run-nothing"
BINARY="run-nothing"

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin) OS_NAME="macos" ;;
  Linux)  OS_NAME="linux" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH_NAME="x86_64" ;;
  arm64|aarch64) ARCH_NAME="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

ASSET_NAME="${BINARY}-${OS_NAME}-${ARCH_NAME}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}"

echo "Downloading ${ASSET_NAME} from GitHub Releases..."
TMP_FILE="$(mktemp)"
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"
chmod +x "$TMP_FILE"

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  if [ -d "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
  else
    echo "Installing to /usr/local/bin (may require sudo)..."
    sudo mv "$TMP_FILE" "/usr/local/bin/${BINARY}"
    echo "✔ Installed ${BINARY} to /usr/local/bin/${BINARY}"
    echo "Run 'run-nothing' to start!"
    exit 0
  fi
fi

mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY}"
echo "✔ Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
echo "Run 'run-nothing' to start!"
