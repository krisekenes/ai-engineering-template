#!/bin/sh
set -eu
cd "$(dirname "$0")/../backend"
unformatted=$(gofmt -l cmd internal)
if [ -n "$unformatted" ]; then
  printf 'Run make fmt; unformatted files:\n%s\n' "$unformatted" >&2
  exit 1
fi
go test -race -count=1 ./...
go vet ./...
go build ./...

cd ../examples/repair-loop
unformatted=$(gofmt -l main.go main_test.go)
if [ -n "$unformatted" ]; then
  printf 'Unformatted repair demo files:\n%s\n' "$unformatted" >&2
  exit 1
fi
go test -race -count=1 ./...
go vet ./...
