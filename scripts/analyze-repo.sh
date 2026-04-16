#!/bin/bash
# Usage: analyze-repo.sh <org/repo> <results-base-dir>
#
# Clones a repository, runs the architecture analyzer, and cleans up.
set -euo pipefail

REPO="${1:?Usage: analyze-repo.sh <org/repo> <results-base-dir>}"
RESULTS_BASE="${2:?Usage: analyze-repo.sh <org/repo> <results-base-dir>}"

# Validate repo name to prevent command injection
if [[ ! "$REPO" =~ ^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$ ]]; then
    echo "[!] Invalid repo format: $REPO (expected: org/repo)" >&2
    exit 1
fi

SHORT="${REPO##*/}"
OUTDIR="${RESULTS_BASE}/${SHORT}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ANALYZER_DIR="$(dirname "${SCRIPT_DIR}")"
ANALYZER_BIN="${ANALYZER_DIR}/rhoai-analyzer"

# Build if needed
if [ ! -f "${ANALYZER_BIN}" ]; then
    echo "[*] Building rhoai-analyzer..."
    (cd "${ANALYZER_DIR}" && go build -o rhoai-analyzer ./cmd/rhoai-analyzer/)
fi

mkdir -p "${OUTDIR}"

CLONE_DIR="/tmp/rhoai-analyzer-repos/${SHORT}"

# Clone (shallow)
echo "[*] Cloning ${REPO}..."
rm -rf "${CLONE_DIR}"
git clone --depth 1 "https://github.com/${REPO}.git" "${CLONE_DIR}" 2>/dev/null || {
    echo "[!] Failed to clone ${REPO}" >&2
    exit 1
}

# Full analysis (architecture + code graph)
echo "[*] Analyzing ${SHORT}..."
"${ANALYZER_BIN}" full-analysis "${CLONE_DIR}" --output-dir "${OUTDIR}"

# Cleanup
rm -rf "${CLONE_DIR}"
echo "[*] Done: ${OUTDIR}"
