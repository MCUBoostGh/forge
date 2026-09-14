#!/usr/bin/env bash

set -u

fail() {
    echo "run.sh: $1" >&2
    exit 1
}

trap 'fail "command failed on line $LINENO"' ERR

usage() {
    cat <<'EOF'
Usage: ./run.sh <project_name> [--device <device>]

Examples:
  ./run.sh blink
  ./run.sh blink --device stm32f103r8
  ./run.sh blink --device bluepill

Default device is stm32f103r8 when --device is omitted.
EOF
}

EXAMPLES_DIR="examples"
DEFAULT_DEVICE="stm32f103r8"

PROJECT_NAME=""
DEVICE="$DEFAULT_DEVICE"

while [ "$#" -gt 0 ]; do
    case "$1" in
        -h|--help)
            usage
            exit 0
            ;;
        --device)
            if [ "$#" -lt 2 ]; then
                fail "--device requires a device id (e.g. stm32f103r8)"
            fi
            DEVICE="$2"
            shift 2
            ;;
        -*)
            fail "unknown option: $1 (try --help)"
            ;;
        *)
            if [ -n "$PROJECT_NAME" ]; then
                fail "unexpected argument: $1"
            fi
            PROJECT_NAME="$1"
            shift
            ;;
    esac
done

if [ -z "$PROJECT_NAME" ]; then
    read -r -p "Enter project name: " PROJECT_NAME
fi

if [ -z "$PROJECT_NAME" ]; then
    fail "project name is required"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR" || fail "unable to enter repo root $SCRIPT_DIR"

mkdir -p "$EXAMPLES_DIR" || fail "unable to create $EXAMPLES_DIR"

go build -o forge . || fail "go build failed"
go install . || fail "go install failed"

# Prefer the freshly built binary; fall back to PATH (go install).
if [ -x "$SCRIPT_DIR/forge" ]; then
    FORGE="$SCRIPT_DIR/forge"
elif command -v forge >/dev/null 2>&1; then
    FORGE="$(command -v forge)"
else
    fail "forge binary not found after build/install"
fi

PROJECT_DIR="$EXAMPLES_DIR/$PROJECT_NAME"
if [ -e "$PROJECT_DIR" ]; then
    fail "project already exists: $PROJECT_DIR"
fi

(
    cd "$EXAMPLES_DIR" || fail "unable to enter $EXAMPLES_DIR"
    "$FORGE" new "$PROJECT_NAME" --device "$DEVICE" || fail "forge new failed"
)

(
    cd "$PROJECT_DIR" || fail "unable to enter $PROJECT_DIR"
    "$FORGE" init || fail "forge init failed"
    "$FORGE" build || fail "forge build failed"
)

echo "Project '$PROJECT_NAME' (device: $DEVICE) created in '$PROJECT_DIR'"
