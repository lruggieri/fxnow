#!/bin/bash

if [[ ! "$PATH" =~ $(go env GOPATH) ]]; then
  PATH="$(go env GOPATH)/bin:$PATH"
  export PATH
fi

# adaptive read yaml field
# $1 is yaml field path (should not start with dot)
# $2 is file path
yqr() {
  ypath="$1"
  filepath="$2"
  # major version
  mv=$(yq --version | sed -E 's/.*(([0-9]+)\.([0-9]+)\.([0-9]+))/\2/')
  if [[ $mv -lt 4 ]]; then
    yq r "$filepath" "$ypath"
  else
    yq eval ".$ypath" "$filepath"
  fi
}

log_info() {
  if [ "$CI" == "true" ]; then
    echo "[INFO] $1"
  else
    COLOR='\033[0;36m'
    NC='\033[0m' # No Color
    printf "${COLOR}[INFO]${NC} %s\n" "$1"
  fi
}

log_error() {
  if [ "$CI" == "true" ]; then
    echo "[ERROR] $1"
  else
    COLOR='\033[0;31m'
    NC='\033[0m' # No Color
    printf "${COLOR}[ERROR]${NC} %s\n" "$1"
  fi
}

pushd() {
  command "pushd" "$@" >/dev/null || exit 1
}

popd() {
  command "popd" "$@" >/dev/null || exit 1
}

removeline() {
  sed -e "/$1/d" "$2" >"$2.tmp" && mv "$2.tmp" "$2"
}

start_time() {
  date +%s.%N
}

since_time() {
  start="$1"
  end=$(date +%s.%N)
  echo "$end - $start" | bc
}
