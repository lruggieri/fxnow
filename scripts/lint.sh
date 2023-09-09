#!/bin/bash
set -e

SCRIPTS_DIR=$(cd "$(dirname "$0")" && pwd)
TARGET_DIR="$1"
TOOL_VERSION="v1.52.2"
TOOL_VERSION_NUMBER=${TOOL_VERSION/v/}

# shellcheck source=./_functions.sh
source "$SCRIPTS_DIR/_functions.sh"

if [[ -z "$TARGET_DIR" ]]; then
  log_error "target dir should not empty"
  exit 1
fi

pushd "$TARGET_DIR"

if [[ ! "$PATH" =~ $(go env GOPATH) ]]; then
  PATH="$(go env GOPATH)/bin:$PATH"
  export PATH
fi

install() {
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" "$TOOL_VERSION"
}

if ! command -v "golangci-lint" >/dev/null 2>&1; then
  install
else
  current_version=$(golangci-lint version 2>&1)
  log_info "current_version=$current_version"
  if ! (echo "$current_version" | grep -Eq "$TOOL_VERSION_NUMBER"); then
    log_info "install required golangci-lint verion: $TOOL_VERSION"
    install
  fi
fi

log_info "linting..."

golangci-lint run --config "$SCRIPTS_DIR/../.golangci.yml" --timeout 5m --verbose ./...
