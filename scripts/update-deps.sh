#!/bin/bash
set -euo pipefail

git config --local user.email "action@github.com"
git config --local user.name "GitHub Action"

export UPDATER_REF="${GITHUB_HEAD_REF:-$GITHUB_REF_NAME}"

node scripts/update-deps.js
