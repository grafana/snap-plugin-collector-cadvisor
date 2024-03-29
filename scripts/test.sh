#!/usr/bin/env bash

# Support travis.ci environment matrix:
TEST_TYPE="${TEST_TYPE:-$1}"
UNIT_TEST="${UNIT_TEST:-"gofmt goimports go_vet go_test go_cover"}"
TEST_K8S="${TEST_K8S:-0}"

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

_debug "script directory ${__dir}"
_debug "project directory ${__proj_dir}"
_info "skipping go test in following directories: ${NO_GO_TEST}"

[[ "$TEST_TYPE" =~ ^(small|medium|large|legacy|build)$ ]] || _error "invalid TEST_TYPE (value must be 'small', 'medium', 'large', 'legacy', or 'build' recieved:${TEST_TYPE}"


test_unit() {
  # The script does automatic checking on a Go package and its sub-packages, including:
  # 1. gofmt         (http://golang.org/cmd/gofmt/)
  # 2. goimports     (https://github.com/bradfitz/goimports)
  # 3. golint        (https://github.com/golang/lint)
  # 4. go vet        (http://golang.org/cmd/vet)
  # 5. race detector (http://blog.golang.org/race-detector)
  # 6. go test
  # 7. test coverage (http://blog.golang.org/cover)
  local go_tests
  go_tests=(gofmt goimports golint go_vet go_race go_test go_cover)

  _debug "available unit tests: ${go_tests[*]}"
  _debug "user specified tests: ${UNIT_TEST}"

  ((n_elements=${#go_tests[@]}, max=n_elements - 1))

  for ((i = 0; i <= max; i++)); do
    if [[ "${UNIT_TEST}" =~ (^| )"${go_tests[i]}"( |$) ]]; then
      _info "running ${go_tests[i]}"
      _"${go_tests[i]}"
    else
      _debug "skipping ${go_tests[i]}"
    fi
  done
}

if [[ $TEST_TYPE == "legacy" ]]; then
  UNIT_TEST="go_test go_cover"
  echo "mode: count" > profile.cov
  export TEST_TYPE="unit"
  test_unit
elif [[ $TEST_TYPE == "small" ]]; then
  if [[ -f "${__dir}/small.sh" ]]; then
    . "${__dir}/small.sh"
  else
    echo "mode: count" > profile.cov
    test_unit
  fi
elif [[ $TEST_TYPE == "medium" ]]; then
  if [[ -f "${__dir}/medium.sh" ]]; then
    . "${__dir}/medium.sh"
  else
    UNIT_TEST="go_test go_cover"
    echo "mode: count" > profile.cov
    test_unit
  fi
elif [[ $TEST_TYPE == "large" ]]; then
  if [[ "${TEST_K8S}" != "0" && -f "$__dir/large_k8s.sh" ]]; then
    . "${__dir}/large_k8s.sh"
  elif [[ -f "${__dir}/large_compose.sh" ]]; then
    . "${__dir}/large_compose.sh"
  else
    . "${__dir}/large.sh"
  fi
elif [[ $TEST_TYPE == "build" ]]; then
  "${__dir}/build.sh"
fi