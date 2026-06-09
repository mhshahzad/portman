#!/bin/bash

set -e

REPO="mhshahzad/portman"
BINARY_NAME="portman"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${OS}" in
    linux*)  OS='linux';;
    darwin*) OS='darwin';;
    *)       echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)  ARCH='amd64';;
    arm64|aarch64) ARCH='arm64';;
    *)       echo "Unsupported Architecture: ${ARCH}"; exit 1;;
esac

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"

echo "Downloading ${BINARY_NAME} for ${OS}/${ARCH}..."
if command -v curl >/dev/null 2>&1; then
    curl -L "${DOWNLOAD_URL}" -o "${BINARY_NAME}"
elif command -v wget >/dev/null 2>&1; then
    wget -O "${BINARY_NAME}" "${DOWNLOAD_URL}"
else
    echo "Error: curl or wget is required."
    exit 1
fi

chmod +x "${BINARY_NAME}"

echo "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
if [ -w "${INSTALL_DIR}" ]; then
    mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    sudo mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "Successfully installed ${BINARY_NAME} $(${INSTALL_DIR}/${BINARY_NAME} version | cut -d' ' -f2)"
