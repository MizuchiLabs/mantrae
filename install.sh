#!/usr/bin/env sh
set -eu

REPO="https://github.com/mizuchilabs/mantrae/releases/download"
REPO_API="https://api.github.com/repos/mizuchilabs/mantrae/releases"
BINARY="mantrae"

main() {
   if [ "${1:-}" = "uninstall" ]; then
      uninstall
      exit 0
   fi

   platform="$(uname -s)"
   arch="$(uname -m)"
   tempdir="$(mktemp -d -t "${BINARY}-XXXXXX")"
   latest=$(curl -fsSL "${REPO_API}/latest" | grep -o '"tag_name":.*' | cut -d '"' -f 4)

   case "$platform" in
   Darwin) platform="darwin" ;;
   Linux) platform="linux" ;;
   MINGW* | MSYS* | CYGWIN*) platform="windows" ;;
   *) echo "Unsupported platform: $platform" && exit 1 ;;
   esac

   case "$arch" in
   arm64* | aarch64*) arch="arm64" ;;
   x86_64* | amd64*) arch="amd64" ;;
   *) echo "Unsupported architecture: $arch" && exit 1 ;;
   esac

   filename="${BINARY}_${platform}_${arch}"
   [ "$platform" = "windows" ] && filename="${filename}.exe"
   url="${REPO}/${latest}/${filename}"

   echo "Downloading ${BINARY} ${latest}..."
   if ! curl -fsSL "$url" -o "$tempdir/$filename"; then
      echo "Download failed"
      rm -rf "$tempdir"
      exit 1
   fi

   mkdir -p "$HOME/.local/bin"
   mv "$tempdir/$filename" "$HOME/.local/bin/$BINARY"
   chmod +x "$HOME/.local/bin/$BINARY"
   rm -rf "$tempdir"

   if echo "$PATH" | grep -q "$HOME/.local/bin"; then
      echo "${BINARY} has been installed. Run with: ${BINARY}"
   else
      echo "Add ~/.local/bin to your PATH:"
      echo "  export PATH=\$HOME/.local/bin:\$PATH"
   fi
}

uninstall() {
   bin_path="$HOME/.local/bin/$BINARY"
   if [ -f "$bin_path" ]; then
      rm -f "$bin_path"
      echo "${BINARY} has been removed."
   else
      echo "${BINARY} is not installed."
   fi
}

main "$@"
