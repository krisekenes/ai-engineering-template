#!/bin/sh
set -eu
for tool in go git make; do
  command -v "$tool" >/dev/null 2>&1 || { echo "Missing prerequisite: $tool" >&2; exit 1; }
done
go version
git --version
printf '%s\n' 'Required: Go 1.25+. Run make check to verify toolchain compatibility.' 'No credentials or external services required.'
