#!/usr/bin/env sh
set -eu

# GitHub releases
REPO="https://github.com/mizuchilabs/mantrae/releases/download"

# GitHub release API
REPO_API="https://api.github.com/repos/mizuchilabs/mantrae/releases"

binary="mantrae"

# Downloads the latest release and moves it into ~/.local/bin
main() {
   case "${1:-}" in
   agent) binary="mantrae-agent" ;;
   uninstall)
      uninstall "${2:-}"
      exit 0
      ;;
   esac

   platform="$(uname -s)"
   arch="$(uname -m)"
   temp="$(mktemp -t "$binary-XXXXXX")"
   latest=$(curl -fsSL "${REPO_API}/latest" | grep -o '"tag_name":.*' | cut -d '"' -f 4)

   case "$platform" in
   Darwin) platform="macos" ;;
   Linux) platform="linux" ;;
   *)
      echo "Unsupported platform: $platform"
      exit 1
      ;;
   esac

   case "$arch" in
   arm64* | aarch64*) arch="arm64" ;;
   x86_64* | amd64*) arch="amd64" ;;
   i?86*) arch="386" ;;
   *)
      echo "Unsupported architecture: $arch"
      exit 1
      ;;
   esac

   url="${REPO}/${latest}/${binary}_${platform}_${arch}"

   echo "Downloading $binary from $url"
   download "$url" "$temp"

   # Ensure the file is not empty
   if [ ! -s "$temp" ]; then
      echo "Download failed: file is empty"
      exit 1
   fi

   install_binary "$temp" "$binary"
   post_install "$binary"
}

download() {
   url="$1"
   dest="$2"
   if command -v curl >/dev/null 2>&1; then
      curl -fsSL "$url" -o "$dest"
   elif command -v wget >/dev/null 2>&1; then
      wget -qO "$dest" "$url"
   else
      echo "Please install 'curl' or 'wget' to proceed"
      exit 1
   fi
}

install_binary() {
   binary_path="$HOME/.local/bin"
   mkdir -p "$binary_path"
   mv "$1" "$binary_path/$binary"
   chmod +x "$binary_path/$binary"
}

post_install() {
   if echo "$PATH" | grep -q "$HOME/.local/bin"; then
      echo "$binary has been installed. Run with $binary"
   else
      echo "Add ~/.local/bin to your PATH to use $binary:"
      case "$SHELL" in
      *zsh) echo "  echo 'export PATH=\$HOME/.local/bin:\$PATH' >> ~/.zshrc && source ~/.zshrc" ;;
      *fish) echo "  fish_add_path -U $HOME/.local/bin" ;;
      *) echo "  echo 'export PATH=\$HOME/.local/bin:\$PATH' >> ~/.bashrc && source ~/.bashrc" ;;
      esac
   fi
}

uninstall() {
   target="${1:-mantrae}"
   case "$target" in
   mantrae | agent)
      binary="mantrae"
      [ "$target" = "agent" ] && binary="mantrae-agent"
      bin_path="$HOME/.local/bin/$binary"
      if [ -f "$bin_path" ]; then
         echo "Removing $binary..."
         rm -f "$bin_path"
         echo "$binary has been removed."
      else
         echo "$binary is not installed."
      fi
      ;;
   *)
      echo "Usage: uninstall [mantrae|agent]"
      exit 1
      ;;
   esac
}

main "$@"
