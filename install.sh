#!/bin/sh
set -e

# Greenlight installer
# Usage: curl -fsSL https://raw.githubusercontent.com/atlantic-blue/greenlight/main/install.sh | sh

REPO="atlantic-blue/greenlight"
INSTALL_DIR="/usr/local/bin"
BINARY="greenlight"

main() {
    detect_platform
    get_latest_version
    download_and_install
    verify_install
}

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux) ;;
        darwin) ;;
        *)
            echo "Error: unsupported OS: $OS"
            echo "Greenlight supports linux and darwin (macOS)."
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *)
            echo "Error: unsupported architecture: $ARCH"
            echo "Greenlight supports amd64 and arm64."
            exit 1
            ;;
    esac

    echo "Platform: ${OS}/${ARCH}"
}

get_latest_version() {
    echo "Fetching latest version..."

    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' \
        | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        echo "Error: could not determine latest version."
        echo "Check https://github.com/${REPO}/releases"
        exit 1
    fi

    # Strip leading 'v' for the archive filename
    VERSION_NUMBER=$(echo "$VERSION" | sed 's/^v//')

    echo "Latest version: ${VERSION}"
}

download_and_install() {
    ARCHIVE="${BINARY}_${VERSION_NUMBER}_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
    TMP_DIR=$(mktemp -d)

    echo "Downloading ${URL}..."

    if ! curl -fsSL -o "${TMP_DIR}/${ARCHIVE}" "$URL"; then
        echo "Error: download failed."
        echo "Check that a release exists for ${OS}/${ARCH} at:"
        echo "  https://github.com/${REPO}/releases/tag/${VERSION}"
        rm -rf "$TMP_DIR"
        exit 1
    fi

    echo "Extracting..."
    tar xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"

    if [ ! -f "${TMP_DIR}/${BINARY}" ]; then
        echo "Error: binary not found in archive."
        rm -rf "$TMP_DIR"
        exit 1
    fi

    echo "Installing to ${INSTALL_DIR}/${BINARY}..."

    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    else
        echo "(requires sudo)"
        sudo mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    fi

    chmod +x "${INSTALL_DIR}/${BINARY}"
    rm -rf "$TMP_DIR"
}

verify_install() {
    if command -v "$BINARY" > /dev/null 2>&1; then
        echo ""
        echo "Greenlight installed successfully!"
        echo ""
        "$BINARY" version
        echo ""
        echo "Get started:"
        echo "  greenlight install --global"
        echo "  greenlight help"
    else
        echo ""
        echo "Installed to ${INSTALL_DIR}/${BINARY} but it's not in your PATH."
        echo "Add ${INSTALL_DIR} to your PATH, or move the binary:"
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

main
