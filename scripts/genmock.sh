#!/bin/bash
set -e

SCRIPTS_DIR=$(cd "$(dirname "$0")" && pwd)
TARGET_DIR="$1"
MOCKERY_VERSION="v2.20.2"
TOOL_VERSION_NUMBER=${MOCKERY_VERSION/v/}
GOBIN=$(go env GOBIN)
EXCLUDES="$2"

if [[ "$GOBIN" == "" ]]; then
  GOBIN="$(go env GOPATH)/bin"
fi

source "$SCRIPTS_DIR/_functions.sh"

if [[ -z "$TARGET_DIR" ]]; then
  log_error "target dir should not empty"
  exit 1
fi

pushd "$TARGET_DIR"

install_by_download() {
  os=$(uname)
  arch=$(uname -m)
  name="mockery_${TOOL_VERSION_NUMBER}_${os}_${arch}.tar.gz"
  url="https://github.com/vektra/mockery/releases/download/$MOCKERY_VERSION/$name"

  output_path="$HOME/Downloads/$name"
  # make sure output dir exist
  mkdir -p "$(dirname "$output_path")"

  log_info "Downloading: $url"

  curl -Lo "$output_path" "$url"

  if [[ -e "$GOBIN/mockery" ]]; then
    rm "$GOBIN/mockery"
  fi

  tar -C "$GOBIN" -xzf "$output_path"
}

if ! command -v "$GOBIN/mockery" >/dev/null 2>&1; then
  install_by_download
else
  current_version=$($GOBIN/mockery --version 2>/dev/null)
  log_info "current_version=$current_version"
  if ! (echo "$current_version" | grep -Eq "$MOCKERY_VERSION"); then
    log_info "install required mockery verion: $MOCKERY_VERSION"
    install_by_download
  fi
fi

rm -rf ./mock

$GOBIN/mockery --all \
  --recursive \
  --output mock \
  --case underscore \
  --keeptree \
  --packageprefix "mock" \
  --with-expecter

for a in $EXCLUDES; do
  f="./mock/$a"

  if [[ -d "$f" ]]; then
    echo "removing: $f"
    rm -rf "$f" || true
  elif [[ -f "$f" ]]; then
    echo "removing: $f"
    rm "$f"
  fi

done
