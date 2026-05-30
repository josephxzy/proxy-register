#!/usr/bin/env bash
set -euo pipefail

# Usage: ./build.sh [version]
# Example: ./build.sh v1.0.0

VERSION=${1:-"v1.0.0"}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$SCRIPT_DIR/src"
FRONTEND_DIR="$SRC_DIR/frontend"
OUTPUT_DIR="$SCRIPT_DIR/build"

mkdir -p "$OUTPUT_DIR"
cd "$SCRIPT_DIR"

echo "============================================"
echo " GHProxy Registry Build Script"
echo " Version: $VERSION"
echo "============================================"
echo ""

echo "[1/3] Building frontend..."
cd "$FRONTEND_DIR"
if [ ! -d "node_modules" ]; then
  npm install
fi
npm run build
cd "$SCRIPT_DIR"
echo "  Frontend built -> src/public/"
echo ""

BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

PLATFORMS=(
  "linux/amd64"
  "windows/amd64"
)

echo "[2/3] Building Linux amd64..."
cd "$SRC_DIR"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/github-proxy-registry-linux-amd64" .
echo "  Done -> build/github-proxy-registry-linux-amd64"

echo "[3/3] Building Windows amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/github-proxy-registry-windows-amd64.exe" .
echo "  Done -> build/github-proxy-registry-windows-amd64.exe"

echo ""

cp "$SRC_DIR/config.toml" "$OUTPUT_DIR/"
echo "Config copied to build/"
echo ""

echo "============================================"
echo " Build complete!"
echo " Output: $OUTPUT_DIR/"
echo ""
echo "Files:"
ls -lh "$OUTPUT_DIR"/github-proxy-registry-*
echo ""
echo "============================================"
