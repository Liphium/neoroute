#!/bin/sh
set -eu
project_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_dir=$(CDPATH= cd -- "$project_dir/../../../.." && pwd)
cd "$repo_dir"
go run ./cmd/neogen/tests/generate -language go -output "$project_dir"
cd "$project_dir"
go mod tidy
go run github.com/tinylib/msgp@latest -file models.go -tests=false
go mod tidy
go test -v ./...
