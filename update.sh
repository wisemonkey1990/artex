#!/usr/bin/env bash
# ARTEX 更新脚本：① Docker 更新（拉取新镜像并重新编排）  ② 本地编译更新（重新构建二进制）
# 与 install.sh 配套。install 负责首次安装，update 负责升级到新版本。
# 不需要手动执行数据库迁移。artex 每次启动都会幂等地重新执行 schema.sql（包含 ADD COLUMN/CREATE
# INDEX IF NOT EXISTS），因此"重启即迁移"。数据（pgdata 卷、./data、./skills）不受影响。
set -euo pipefail
cd "$(cd "$(dirname "$0")" && pwd)"

info(){ printf '\033[36m[*]\033[0m %s\n' "$*"; }
ok(){   printf '\033[32m[+]\033[0m %s\n' "$*"; }
warn(){ printf '\033[33m[!]\033[0m %s\n' "$*"; }
die(){  printf '\033[31m[x]\033[0m %s\n' "$*" >&2; exit 1; }
ask(){  local p="$1" d="${2:-}" a; read -rp "$p${d:+ [$d]}: " a; echo "${a:-$d}"; }

# ── 可选：把仓库同步到最新代码（compose·脚本·本地编译源码都会一起更新） ───────
sync_repo(){
  [ -d .git ] && command -v git >/dev/null 2>&1 || { warn "不是 git 工作副本，跳过 git pull"; return; }
  [ "$(ask '是否拉取最新代码（git pull --ff-only）？(y/n)' y)" = y ] || return
  if ! git pull --ff-only; then
    warn "git pull 无法快进（存在本地改动或分支已分叉）。请自行处理后再试，本次将直接使用当前代码"
  fi
}

# ── ① Docker 更新 ───────────────────────────────
update_docker(){
  command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1 \
    || die "未找到 docker / docker compose，请先用 ./install.sh 完成安装部署"
  [ -f .env ] || die "未找到 .env，请先用 ./install.sh 完成首次部署"

  # 可选：升级到指定的版本 tag（留空则沿用 .env 中的 ARTEX_TAG，默认为 latest）
  local tag; tag="$(ask '目标镜像 tag（回车使用 .env / latest）' '')"
  if [ -n "$tag" ]; then
    if grep -q '^ARTEX_TAG=' .env; then
      sed -i.bak "s|^ARTEX_TAG=.*|ARTEX_TAG=${tag}|" .env && rm -f .env.bak
    else
      printf '\nARTEX_TAG=%s\n' "$tag" >> .env
    fi
    ok "已将 ARTEX_TAG 设置为 ${tag}"
  fi

  # 只更新 artex。postgres 固定用 16-alpine，不需要跟着升级（拉取只是浪费带宽，
  # 且大版本变化还有兼容性风险）。artex 声明了 depends_on postgres，所以带服务名
  # up 时若 pg 未运行会自动启动，若已运行则保持不变、不重新编排。
  info "正在拉取新镜像（仅 artex）…"
  docker compose pull artex
  info "正在重新编排并启动（artex 重启时会自动迁移 schema）…"
  docker compose up -d artex
  ok "更新完成 → http://localhost:8787"
  warn "刚拉取的镜像是上游官方 autumn27/artex 镜像，不包含本 fork 的改动，请用\"2) 本地更新（go 重新编译）\"来构建"
  info "查看日志：docker compose logs -f artex"
  info "清理旧镜像（可选）：docker image prune -f"
}

# ── ② 本地编译更新 ──────────────────────────────
update_local(){
  command -v go >/dev/null 2>&1 || die "未找到 Go（>=1.26）：https://go.dev/dl/"
  [ -f config.json ] || warn "未找到 config.json。若是首次部署，请使用 ./install.sh"
  ok "Go: $(go version)"

  if command -v npm >/dev/null 2>&1; then
    info "正在重新构建前端静态产物…"
    ( cd web && npm ci && npm run build:static )
    rm -rf server/webui/dist && cp -r web/out server/webui/dist
    info "正在重新编译内嵌前端的单一二进制…"
    CGO_ENABLED=0 go build -tags embedui -trimpath -o artex ./cmd/artex
  else
    warn "未找到 npm，仅编译不内嵌前端的后端（前端需另行用 npm run dev 运行）"
    CGO_ENABLED=0 go build -o artex ./cmd/artex
  fi
  ok "编译完成 → ./artex"
  warn "请重启运行中的 artex 进程以应用更改（重启时会自动迁移 schema）"
}

echo "=============================="
echo "  ARTEX 更新"
echo "  1) Docker 更新（拉取新镜像并重新编排）"
echo "  2) 本地更新（go 重新编译）"
echo "=============================="
case "$(ask '选择' 1)" in
  1) sync_repo; update_docker ;;
  2) sync_repo; update_local ;;
  *) die "选择无效" ;;
esac
