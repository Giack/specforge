#!/usr/bin/env bash
# install.sh — SpecForge binary installer
# Usage: curl -sSL https://raw.githubusercontent.com/Giack/specforge/main/install.sh | bash
# Or:    ./install.sh [--prefix /usr/local]
set -euo pipefail

REPO="Giack/specforge"
BINARY="specforge"
DEFAULT_PREFIX="${HOME}/.local"

# ── helpers ───────────────────────────────────────────────────────────────────

info()  { printf '\033[0;34m[specforge]\033[0m %s\n' "$*"; }
ok()    { printf '\033[0;32m[specforge]\033[0m %s\n' "$*"; }
die()   { printf '\033[0;31m[specforge] error:\033[0m %s\n' "$*" >&2; exit 1; }

# ── plugin installer ──────────────────────────────────────────────────────────

install_plugin() {
  if ! command -v claude &>/dev/null; then
    info "Claude Code CLI not found — skipping plugin install."
    info "To install manually: claude plugin marketplace add Giack/specforge && claude plugin install specforge"
    return
  fi

  info "Installing Claude Code plugin via claude CLI..."

  # Remove stale marketplace entry if present, then re-add from GitHub
  claude plugin marketplace remove specforge 2>/dev/null || true
  claude plugin marketplace add Giack/specforge \
    || die "Failed to add specforge marketplace"

  # Remove stale plugin if present, then install fresh
  claude plugin uninstall specforge --scope user 2>/dev/null || true
  claude plugin install specforge@specforge --scope user \
    || die "Failed to install specforge plugin"

  ok "Plugin installed. Restart Claude Code to activate /specforge:map"
}

# ── detect platform ───────────────────────────────────────────────────────────

detect_os() {
  case "$(uname -s)" in
    Darwin) echo "darwin" ;;
    Linux)  echo "linux"  ;;
    *)      die "Unsupported OS: $(uname -s). Please build from source." ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64 | amd64)  echo "amd64"  ;;
    arm64  | aarch64) echo "arm64" ;;
    *)                die "Unsupported arch: $(uname -m). Please build from source." ;;
  esac
}

# ── parse args ────────────────────────────────────────────────────────────────

PREFIX="${DEFAULT_PREFIX}"
VERSION="latest"
INSTALL_PLUGIN=true

while [[ $# -gt 0 ]]; do
  case "$1" in
    --prefix) PREFIX="$2"; shift 2 ;;
    --version) VERSION="$2"; shift 2 ;;
    --no-plugin) INSTALL_PLUGIN=false; shift ;;
    -h|--help)
      echo "Usage: install.sh [--prefix DIR] [--version TAG] [--no-plugin]"
      echo ""
      echo "  --prefix DIR     Install to DIR/bin (default: ~/.local)"
      echo "  --version TAG    Install specific release tag (default: latest)"
      echo "  --no-plugin      Skip Claude Code plugin installation"
      exit 0
      ;;
    *) die "Unknown argument: $1" ;;
  esac
done

BIN_DIR="${PREFIX}/bin"

# ── resolve version ───────────────────────────────────────────────────────────

if [[ "${VERSION}" == "latest" ]]; then
  info "Resolving latest release..."
  if command -v curl &>/dev/null; then
    VERSION="$(curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" \
      | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
  elif command -v wget &>/dev/null; then
    VERSION="$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" \
      | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
  else
    die "Neither curl nor wget found. Please install one of them."
  fi
  [[ -z "${VERSION}" ]] && die "Could not resolve latest release version."
fi

OS="$(detect_os)"
ARCH="$(detect_arch)"
ASSET="${BINARY}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"

info "Installing specforge ${VERSION} (${OS}/${ARCH})"
info "Destination: ${BIN_DIR}/${BINARY}"

# ── download ──────────────────────────────────────────────────────────────────

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

TMP_BIN="${TMP_DIR}/${BINARY}"

info "Downloading ${DOWNLOAD_URL} ..."
if command -v curl &>/dev/null; then
  curl -sSfL --output "${TMP_BIN}" "${DOWNLOAD_URL}" \
    || die "Download failed. Check that release ${VERSION} has asset ${ASSET}."
else
  wget -qO "${TMP_BIN}" "${DOWNLOAD_URL}" \
    || die "Download failed. Check that release ${VERSION} has asset ${ASSET}."
fi

chmod +x "${TMP_BIN}"

# ── verify binary ─────────────────────────────────────────────────────────────

info "Verifying binary..."
"${TMP_BIN}" --version 2>/dev/null | head -1 || true  # non-fatal if stub

# ── install ───────────────────────────────────────────────────────────────────

mkdir -p "${BIN_DIR}"
mv "${TMP_BIN}" "${BIN_DIR}/${BINARY}"

ok "specforge ${VERSION} installed to ${BIN_DIR}/${BINARY}"

if [[ "${INSTALL_PLUGIN}" == "true" ]]; then
  install_plugin "${VERSION}"
fi

# ── PATH hint ─────────────────────────────────────────────────────────────────

if ! command -v specforge &>/dev/null; then
  echo ""
  info "Add specforge to your PATH:"
  echo ""

  SHELL_NAME="$(basename "${SHELL:-bash}")"
  case "${SHELL_NAME}" in
    zsh)  RC="~/.zshrc"  ;;
    fish) RC="~/.config/fish/config.fish" ;;
    *)    RC="~/.bashrc" ;;
  esac

  if [[ "${PREFIX}" == "${DEFAULT_PREFIX}" ]]; then
    echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ${RC}"
    echo "  source ${RC}"
  else
    echo "  echo 'export PATH=\"${BIN_DIR}:\$PATH\"' >> ${RC}"
    echo "  source ${RC}"
  fi
  echo ""
fi

if [[ "${INSTALL_PLUGIN}" == "true" ]]; then
  ok "Done! Run: specforge --help  |  In Claude Code: /specforge:map"
else
  ok "Done! Run: specforge --help"
fi
