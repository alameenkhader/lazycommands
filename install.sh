#!/bin/sh
# Install lazycommands from the latest GitHub release.
# Usage: curl -fsSL https://raw.githubusercontent.com/alameenkhader/lazycommands/main/install.sh | sh
set -eu

REPO="alameenkhader/lazycommands"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

case "$os" in
  darwin|linux) ;;
  *) echo "Unsupported OS: $os (use Homebrew or download a binary from https://github.com/$REPO/releases)"; exit 1 ;;
esac

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "Unsupported architecture: $arch"; exit 1 ;;
esac

asset="lazycommands-$os-$arch"
url="https://github.com/$REPO/releases/latest/download/$asset"

install_dir="/usr/local/bin"
if [ ! -w "$install_dir" ]; then
  install_dir="$HOME/.local/bin"
  mkdir -p "$install_dir"
fi

echo "Downloading $asset..."
tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
curl -fsSL "$url" -o "$tmp"

install -m 755 "$tmp" "$install_dir/lazycommands"
ln -sf "$install_dir/lazycommands" "$install_dir/lc"

echo "Installed lazycommands to $install_dir/lazycommands (alias: lc)"
case ":$PATH:" in
  *":$install_dir:"*) ;;
  *) echo "Note: $install_dir is not in your PATH. Add it with:"
     echo "  export PATH=\"$install_dir:\$PATH\"" ;;
esac

"$install_dir/lazycommands" --version
