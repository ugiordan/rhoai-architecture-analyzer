#!/bin/bash
# Usage: analyze-repo.sh <org/repo> <results-base-dir> [version-label]
#
# Clones a repository, runs the architecture analyzer, and cleans up.
# When version-label is provided, output goes to <results-base>/<repo>/<version>/.
set -euo pipefail

# Security hardening for Go toolchain on untrusted repos
export CGO_ENABLED=0
export GOFLAGS="-mod=readonly"
export GONOSUMCHECK=""
export GONOSUMDB=""
export GOMAXPROCS=2

REPO="${1:?Usage: analyze-repo.sh <org/repo> <results-base-dir> [version-label]}"
RESULTS_BASE="${2:?Usage: analyze-repo.sh <org/repo> <results-base-dir> [version-label]}"
VERSION_LABEL="${3:-}"

# SIGTERM/SIGINT trap: when GitHub Actions timeout kills the process,
# clean up and exit 0 so the job shows a warning, not a red error.
trap 'echo "::warning::${REPO}: terminated (timeout or signal)"; \
      chmod -R u+w "${CLONE_DIR:-/dev/null}" 2>/dev/null || true; \
      rm -rf "${CLONE_DIR:-/dev/null}" 2>/dev/null || true; exit 0' TERM INT

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

# Binary location: CI downloads pre-built binary, local dev builds on demand
if [ -x "${GITHUB_WORKSPACE:-}/arch-analyzer" ]; then
    ANALYZER_BIN="${GITHUB_WORKSPACE}/arch-analyzer"
elif [ -x "${ANALYZER_DIR}/arch-analyzer" ]; then
    ANALYZER_BIN="${ANALYZER_DIR}/arch-analyzer"
else
    echo "[*] Building arch-analyzer..."
    (cd "${ANALYZER_DIR}" && go build -o arch-analyzer ./cmd/arch-analyzer/)
    ANALYZER_BIN="${ANALYZER_DIR}/arch-analyzer"
fi

mkdir -p "${OUTDIR}"

CLONE_BASE="${RUNNER_TEMP:-/tmp}/arch-analyzer-repos"
CLONE_DIR="${CLONE_BASE}/${SHORT}"

export GOTMPDIR="${RUNNER_TEMP:-/tmp}/gotmp"
mkdir -p "${GOTMPDIR}"

# Clone (shallow)
echo "[*] Cloning ${REPO}..."
rm -rf "${CLONE_DIR}"
git clone --depth 1 "https://github.com/${REPO}.git" "${CLONE_DIR}" 2>/dev/null || {
    echo "::warning::Skipping ${REPO}: clone failed (repo may be private or inaccessible)"
    exit 0
}

# Symlink boundary check helper
check_symlinks() {
    local resolved_clone
    resolved_clone=$(readlink -f "${CLONE_DIR}" 2>/dev/null || echo "${CLONE_DIR}")
    local escaped
    escaped=$(find "${CLONE_DIR}" -type l -exec readlink -f {} \; 2>/dev/null | grep -v "^${resolved_clone}" || true)
    if [ -n "$escaped" ]; then
        echo "::warning::Skipping ${REPO}: symlinks escape clone boundary"
        rm -rf "${CLONE_DIR}"
        exit 0
    fi
}

# Check symlinks after clone
check_symlinks

# Download Go dependencies for go/packages analysis (isolated cache)
if [ -f "${CLONE_DIR}/go.mod" ]; then
    echo "[*] Downloading Go modules for ${SHORT}..."
    (
        export GOMODCACHE="${CLONE_DIR}/.gomod-cache"
        export GOCACHE="${CLONE_DIR}/.gobuild-cache"
        export GOFLAGS=""
        cd "${CLONE_DIR}" && timeout 120 go mod download 2>&1
    ) || {
        echo "::warning::go mod download failed for ${REPO}, Go AST extraction will use fallback"
    }
    # Re-check symlinks after go mod download (TOCTOU mitigation)
    check_symlinks
fi

# Resolve aliases from scan-config.yaml (if present)
ALIASES_ARGS=""
SCAN_CONFIG="${ANALYZER_DIR}/scan-config.yaml"
if [ -f "${SCAN_CONFIG}" ] && command -v yq &>/dev/null; then
    # Search all platforms for this repo's aliases
    ALIASES=$(yq -r "
        .platforms[].repo_overrides.\"${SHORT}\".aliases // [] | join(\",\")
    " "${SCAN_CONFIG}" 2>/dev/null | grep -v '^$' | head -1 || true)
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

"${ANALYZER_BIN}" full-analysis -output-dir "${OUTDIR}" ${VERSION_ARGS} ${ALIASES_ARGS} "${CLONE_DIR}" || {
    echo "::warning::Analysis failed for ${REPO} (exit $?), partial results may be available"
}

# Cleanup (chmod needed because Go module cache sets files read-only)
chmod -R u+w "${CLONE_DIR}" 2>/dev/null || true
rm -rf "${CLONE_DIR}"
echo "[*] Done: ${OUTDIR}"
