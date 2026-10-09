#!/usr/bin/env bash
# 说明。
# 说明。
#
# 说明。
set -euo pipefail
cd "$(dirname "$0")"

# 说明。
cleanup() { kill 0 2>/dev/null || true; }
trap cleanup EXIT INT TERM

# 说明。
go run ./cmd/artex -addr :8787 -proxy 127.0.0.1:8788 &

# 说明。
( cd web && npm run dev ) &

echo "[dev] 后端 :8787 / 代理 :8788 / 前端 http://localhost:5173（按 Ctrl-C 退出）"
wait
