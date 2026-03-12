#!/bin/bash
# Usage: analyze-repo.sh <org/repo> <results-base-dir>
#
# Clones a repository, runs the architecture analyzer, and cleans up.
set -euo pipefail

REPO="${1:?Usage: analyze-repo.sh <org/repo> <results-base-dir>}"
RESULTS_BASE="${2:?Usage: analyze-repo.sh <org/repo> <results-base-dir>}"
SHORT="${REPO##*/}"
OUTDIR="${RESULTS_BASE}/${SHORT}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ANALYZER_DIR="$(dirname "${SCRIPT_DIR}")"

mkdir -p "${OUTDIR}"

CLONE_DIR="/tmp/rhoai-analyzer-repos/${SHORT}"

# Clone (shallow)
echo "[*] Cloning ${REPO}..."
rm -rf "${CLONE_DIR}"
git clone --depth 1 "https://github.com/${REPO}.git" "${CLONE_DIR}" 2>/dev/null || {
    echo "[!] Failed to clone ${REPO}" >&2
    exit 1
}

# Extract + render
echo "[*] Analyzing ${SHORT}..."
python3 "${ANALYZER_DIR}/analyze.py" analyze "${CLONE_DIR}" --output-dir "${OUTDIR}"

# Cleanup
rm -rf "${CLONE_DIR}"
echo "[*] Done: ${OUTDIR}"
