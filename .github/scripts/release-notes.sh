#!/usr/bin/env bash
# Prints the section for one version from RELEASE_NOTES.md, without its heading.
# Usage: release-notes.sh <version>   (e.g. 4.3.0)
#
# A section starts at a line "## vX.Y.Z" (optionally followed by a space and more text) and ends at the next
# "## " line. Fails if the section is missing or empty, so a release cannot go out without notes.
set -euo pipefail

VERSION=$1
FILE="$(dirname "${BASH_SOURCE[0]}")/../../RELEASE_NOTES.md"

NOTES=$(awk -v heading="## v${VERSION}" '
  index($0, "## ") == 1 {
    if (found) { exit }
    if ($0 == heading || index($0, heading " ") == 1) { found = 1; next }
  }
  found { print }
' "${FILE}")

# Trim leading and trailing blank lines.
NOTES=$(printf '%s\n' "${NOTES}" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}' | sed '/./,$!d')

if [[ -z "${NOTES//[[:space:]]/}" ]]; then
  echo "ERROR: no release notes for v${VERSION}: add a '## v${VERSION}' section to RELEASE_NOTES.md before tagging." >&2
  exit 1
fi

printf '%s\n' "${NOTES}"
