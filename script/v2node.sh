#!/usr/bin/env bash
# Compatibility entry point for nodes using the original manager path.
set -euo pipefail
work=$(mktemp -d)
trap 'rm -rf -- "$work"' EXIT
curl -fLSs --retry 2 'https://raw.githubusercontent.com/Anikato/v2node/v0.4.7.anikato.v2board.1/script/install.sh' -o "$work/v2node"
bash "$work/v2node" "$@"
