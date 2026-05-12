#!/bin/bash
# Usage: analyze-repo.sh <org/repo> <results-base-dir> [version-label]
#
# Clones a repository, runs the architecture analyzer, and cleans up.
# When version-label is provided, output goes to <results-base>/<repo>/<version>/.
set -euo pipefail

REPO="${1:?Usage: analyze-repo.sh <org/repo> <results-base-dir> [version-label]}"
RESULTS_BASE="${2:?Usage: analyze-repo.sh <org/repo> <results-base-dir> [version-label]}"
VERSION_LABEL="${3:-}"

# Validate repo name to prevent command injection
if [[ ! "$REPO" =~ ^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$ ]]; then
    echo "[!] Invalid repo format: $REPO (expected: org/repo)" >&2
    exit 1
fi

ORG="${REPO%%/*}"
SHORT="${REPO##*/}"
OUTDIR="${RESULTS_BASE}/${ORG}/${SHORT}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ANALYZER_DIR="$(dirname "${SCRIPT_DIR}")"
ANALYZER_BIN="${ANALYZER_DIR}/arch-analyzer"

# Build if needed
if [ ! -f "${ANALYZER_BIN}" ]; then
    echo "[*] Building arch-analyzer..."
    (cd "${ANALYZER_DIR}" && go build -o arch-analyzer ./cmd/arch-analyzer/)
fi

mkdir -p "${OUTDIR}"

CLONE_DIR="/tmp/arch-analyzer-repos/${SHORT}"

# Clone (shallow)
echo "[*] Cloning ${REPO}..."
rm -rf "${CLONE_DIR}"
git clone --depth 1 "https://github.com/${REPO}.git" "${CLONE_DIR}" 2>/dev/null || {
    echo "::warning::Skipping ${REPO}: clone failed (repo may be private or inaccessible)"
    exit 0
}

# Resolve aliases from scan-config.yaml (if present)
ALIASES_ARGS=""
SCAN_CONFIG="${ANALYZER_DIR}/scan-config.yaml"
if [ -f "${SCAN_CONFIG}" ] && command -v yq &>/dev/null; then
    # Search all platforms for this repo's aliases
    ALIASES=$(yq -r "
        .platforms[].repo_overrides.\"${SHORT}\".aliases // [] | join(\",\")
    " "${SCAN_CONFIG}" 2>/dev/null | grep -v '^$' | head -1)
    if [ -n "${ALIASES}" ]; then
        ALIASES_ARGS="-aliases ${ALIASES}"
    fi
fi

# Full analysis (architecture + code graph)
echo "[*] Analyzing ${SHORT}..."
VERSION_ARGS=""
if [ -n "${VERSION_LABEL}" ]; then
    VERSION_ARGS="-version ${VERSION_LABEL}"
fi
"${ANALYZER_BIN}" full-analysis -output-dir "${OUTDIR}" ${VERSION_ARGS} ${ALIASES_ARGS} "${CLONE_DIR}"

# Cleanup
rm -rf "${CLONE_DIR}"
echo "[*] Done: ${OUTDIR}"
