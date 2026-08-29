#!/bin/bash
set -euo pipefail

REF="${1:?usage: release-pipe.sh <ref>}"

git config --local user.email "action@github.com"
git config --local user.name "GitHub Action"

export GH_TOKEN="${GH_TOKEN:?GH_TOKEN not set}"

node scripts/get-tree.js 1 "$REF" # transporters
node scripts/get-tree.js 2 "$REF" # neodebug
node scripts/get-tree.js 3 "$REF" # chat_server example
