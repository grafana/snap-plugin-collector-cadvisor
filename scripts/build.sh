#!/bin/bash

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

plugin_name=${__proj_dir##*/}
build_dir="${__proj_dir}/build"
go_build=(go build -ldflags "-w")

_info "project path: ${__proj_dir}"
_info "plugin name: ${plugin_name}"

# rebuild binaries:
_debug "removing: ${build_dir:?}/*"
rm -rf "${build_dir:?}/"*

_info "building plugin: ${plugin_name} for Linux"
export GOOS=linux
export GOARCH=amd64
mkdir -p "${build_dir}/${GOOS}/x86_64"
"${go_build[@]}" -o "${build_dir}/${GOOS}/x86_64/${plugin_name}" . || exit 1

if [[ "$(uname -s)" == "Darwin" ]]; then
_info "building plugin: ${plugin_name} for Darwin"
  export GOOS=darwin
  export GOARCH=amd64
  mkdir -p "${build_dir}/${GOOS}/x86_64"
  "${go_build[@]}" -o "${build_dir}/${GOOS}/x86_64/${plugin_name}" . || exit 1
fi
