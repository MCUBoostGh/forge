#!/usr/bin/env bash

set -u

fail() {
    echo "Error: $1" >&2
    exit 1
}

trap 'fail "command failed on line $LINENO"' ERR

if [ "$#" -gt 0 ]; then
    PROJECT_NAME="$1"
    shift
else
    read -p "Enter project name: " PROJECT_NAME
fi

if [ -z "$PROJECT_NAME" ]; then
    fail "Project name is required."
fi

GO_PROJECT_DIR="examples/"
FORGE_ARGS=("$@")

mkdir -p "$(dirname "$GO_PROJECT_DIR")" || fail "Unable to create examples directory"

go build || fail "go build failed"
go install || fail "go install failed"
cd "$GO_PROJECT_DIR" || fail "Unable to enter $GO_PROJECT_DIR"
forge new "$PROJECT_NAME" "${FORGE_ARGS[@]}" || fail "forge new failed for $GO_PROJECT_DIR"
cd "$PROJECT_NAME" || fail "Unable to enter $GO_PROJECT_DIR"
forge init || fail "forge init failed"
forge build || fail "forge build failed"

echo "Project '$PROJECT_NAME' created successfully in '$GO_PROJECT_DIR'"
