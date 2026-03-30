#!/usr/bin/env bash
set -euo pipefail

# Installs local development prerequisites on Ubuntu:
# - Go (includes gofmt)
# - errcheck
# - goreleaser
# - OpenTofu (tofu) via apt repository

GO_VERSION="${GO_VERSION:-1.23.0}"

log() {
  printf '[devinstaller] %s\n' "$*"
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

as_root() {
  if [[ "${EUID}" -eq 0 ]]; then
    "$@"
  else
    sudo "$@"
  fi
}

detect_arch() {
  case "$(uname -m)" in
    x86_64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *)
      echo "Unsupported architecture: $(uname -m)" >&2
      exit 1
      ;;
  esac
}

ensure_ubuntu() {
  if [[ ! -f /etc/os-release ]]; then
    echo "Cannot detect OS (missing /etc/os-release)." >&2
    exit 1
  fi

  # shellcheck source=/dev/null
  . /etc/os-release
  if [[ "${ID:-}" != "ubuntu" ]]; then
    echo "This installer is intended for Ubuntu. Detected ID=${ID:-unknown}." >&2
    exit 1
  fi
}

install_go() {
  local arch tgz url tmp
  arch="$(detect_arch)"
  tgz="go${GO_VERSION}.linux-${arch}.tar.gz"
  url="https://go.dev/dl/${tgz}"
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' RETURN

  if command -v go >/dev/null 2>&1; then
    if go version | grep -q "go${GO_VERSION}"; then
      log "Go ${GO_VERSION} already installed."
      return
    fi
    log "Replacing existing Go installation."
  fi

  log "Downloading ${url}"
  curl -fsSL "$url" -o "${tmp}/${tgz}"

  log "Installing Go ${GO_VERSION} to /usr/local/go"
  as_root rm -rf /usr/local/go
  as_root tar -C /usr/local -xzf "${tmp}/${tgz}"

  log "Configuring /etc/profile.d/go-path.sh"
  as_root tee /etc/profile.d/go-path.sh >/dev/null <<'EOF'
export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"
EOF

  export PATH="/usr/local/go/bin:${HOME}/go/bin:${PATH}"
}

install_go_tools() {
  log "Installing errcheck and goreleaser"
  as_root env PATH="/usr/local/go/bin:${PATH}" GOBIN=/usr/local/bin \
    /usr/local/go/bin/go install github.com/kisielk/errcheck@latest
  as_root env PATH="/usr/local/go/bin:${PATH}" GOBIN=/usr/local/bin \
    /usr/local/go/bin/go install github.com/goreleaser/goreleaser/v2@latest
}

install_tofu() {
  log "Installing OpenTofu apt repository and package"
  as_root apt-get update
  as_root apt-get install -y ca-certificates curl gpg

  as_root install -d -m 0755 /etc/apt/keyrings
  curl -fsSL https://packages.opentofu.org/opentofu/tofu/gpgkey \
    | as_root gpg --dearmor -o /etc/apt/keyrings/opentofu.gpg
  as_root chmod a+r /etc/apt/keyrings/opentofu.gpg

  echo "deb [signed-by=/etc/apt/keyrings/opentofu.gpg] https://packages.opentofu.org/opentofu/tofu/any/ any main" \
    | as_root tee /etc/apt/sources.list.d/opentofu.list >/dev/null

  as_root apt-get update
  as_root apt-get install -y tofu
}

print_versions() {
  log "Installed versions:"
  go version || true
  gofmt -h >/dev/null 2>&1 && echo "gofmt OK" || true
  errcheck -version || true
  goreleaser --version | head -n1 || true
  tofu version || true
}

main() {
  require_cmd curl
  require_cmd tar
  require_cmd sudo
  ensure_ubuntu

  install_go
  install_go_tools
  install_tofu
  print_versions
}

main "$@"
