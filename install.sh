#!/usr/bin/env sh
set -eu

repo="AndroidPoet/shipkit"
bin_name="shipkit"
install_dir="${INSTALL_DIR:-/usr/local/bin}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$arch" in
  x86_64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac

version="${VERSION:-latest}"
if [ "$version" = "latest" ]; then
  version="$(curl -fsSL "https://api.github.com/repos/$repo/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1)"
fi

archive="${bin_name}_${version}_${os}_${arch}.tar.gz"
url="https://github.com/$repo/releases/download/$version/$archive"
tmp="$(mktemp -d)"

cleanup() {
  rm -rf "$tmp"
}
trap cleanup EXIT

echo "Downloading $url"
curl -fsSL "$url" -o "$tmp/$archive"
tar -xzf "$tmp/$archive" -C "$tmp"

mkdir -p "$install_dir"
cp "$tmp/$bin_name" "$install_dir/$bin_name"
chmod +x "$install_dir/$bin_name"

echo "Installed $bin_name to $install_dir/$bin_name"
