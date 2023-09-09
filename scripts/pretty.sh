#!/bin/bash

set -e

SCRIPTS_DIR=$(cd "$(dirname "$0")" && pwd)
TARGET_DIR="$1"
TOOL_VERSION="v0.3.1"
MODULE=$(go list -m)

# shellcheck source=./_functions.sh
source "$SCRIPTS_DIR/_functions.sh"

if [[ -z "$TARGET_DIR" ]]; then
  log_error "target dir should not empty"
  exit 1
fi

pushd "$TARGET_DIR"

install_gofumpt() {
  go install mvdan.cc/gofumpt@$TOOL_VERSION
}

if ! command -v "gofumpt" >/dev/null 2>&1; then
  install_gofumpt
else
  current_version=$(gofumpt --version 2>/dev/null)
  log_info "current_version=$current_version"
  if ! (echo "$current_version" | grep -Eq "$TOOL_VERSION"); then
    log_info "install required gofumpt verion: $TOOL_VERSION"
    install_gofumpt
  fi
fi

if ! command -v goimports >/dev/null 2>&1; then
  go install golang.org/x/tools/cmd/goimports@latest
fi

log_info "Make code pretty...."
log_info "Below is list of files are changed"

find -H . -name '*.go' -not -path './mock/*' |
  while read -r srcFile; do
    gofumpt -l -w "$srcFile"
    goimports -w -local "$MODULE" "$srcFile"
  done

log_info "--------------------------------"
