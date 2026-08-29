#!/bin/bash

PROJECT="${1:-b}"
if [ ! -d "$PROJECT" ]; then
  echo "error: project directory '$PROJECT' not found"
  exit 1
fi

NEW_VERSION=$(jq -r .\"$PROJECT\" .release-please-manifest.json)
if [ "$NEW_VERSION" = "null" ] || [ -z "$NEW_VERSION" ]; then
  echo "error: no version found for '$PROJECT' in .release-please-manifest.json"
  exit 1
fi

if gh release view "$PROJECT/v$NEW_VERSION" >/dev/null 2>&1; then
  echo "$PROJECT/v$NEW_VERSION already released, skipping"
  exit 0
fi
# Extract the version section from the project's CHANGELOG as release notes
NOTES=$(awk -v v="$NEW_VERSION" 'BEGIN{p=0} /^## \[?'"$NEW_VERSION"'\]?\]?/{p=1; next} p && /^## /{exit} p{print}' "$PROJECT/CHANGELOG.md" | sed '/./,$!d')
git tag "$PROJECT/v$NEW_VERSION"
git push origin "$PROJECT/v$NEW_VERSION"
gh release create "$PROJECT/v$NEW_VERSION" --title "$PROJECT/v$NEW_VERSION" --notes "$NOTES" --latest
