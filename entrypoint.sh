#!/bin/sh

set -e

# Check if command is a valid Mantrae subcommand
if mantrae "$1" --help >/dev/null 2>&1; then
  set -- mantrae "$@"
else
  echo "= '$1' is not a valid subcommand" >&2
  exit 1
fi

exec "$@"
