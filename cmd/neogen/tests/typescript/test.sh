#!/bin/sh
set -eu
project_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_dir=$(CDPATH= cd -- "$project_dir/../../../.." && pwd)
cd "$repo_dir"
go run ./cmd/neogen/tests/generate -language typescript -output "$project_dir/generated"
cd "$project_dir"
npm ci --no-audit --no-fund
npm test
