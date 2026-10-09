#!/usr/bin/env bash
# 说明。
# 说明。
# 说明。
# 说明。
set -euo pipefail
cd "$(cd "$(dirname "$0")" && pwd)"

info(){ printf '\033[36m[*]\033[0m %s\n' "$*"; }
ok(){   printf '\033[32m[+]\033[0m %s\n' "$*"; }
warn(){ printf '\033[33m[!]\033[0m %s\n' "$*"; }
die(){  printf '\033[31m[x]\033[0m %s\n' "$*" >&2; exit 1; }
ask(){  local p="$1" d="${2:-}" a; read -rp "$p${d:+ [$d]}: " a; echo "${a:-$d}"; }

# 说明。
sync_repo(){
  [ -d .git ] && command -v git >/dev/null 2>&1 || { warn "当前目录不是 Git 工作副本，跳过 git pull"; return; }
  [ "$(ask '是否获取最新代码（git pull --ff-only）？(y/n)' y)" = y ] || return
  if ! git pull --ff-only; then
    warn "git pull 无法快进（可能存在本地修改或分支分叉）。请手动处理后重试；本次继续使用当前代码"
  fi
}

# 说明。
update_docker(){
  command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1 \
    || die "未找到 docker / docker compose。请先运行 ./install.sh 完成安装和部署"
  [ -f .env ] || die "找不到 .env。请先运行 ./install.sh 完成首次部署"

  # 可选：升级到指定版本标签；留空则使用 .env 中的 ARTEX_TAG，默认值为 latest。
  local tag; tag="$(ask '目标镜像标签（留空使用 .env / latest）' '')"
  if [ -n "$tag" ]; then
    if grep -q '^ARTEX_TAG=' .env; then
      sed -i.bak "s|^ARTEX_TAG=.*|ARTEX_TAG=${tag}|" .env && rm -f .env.bak
    else
      printf '\nARTEX_TAG=%s\n' "$tag" >> .env
    fi
    ok "已将 ARTEX_TAG 设置为 ${tag}"
  fi

  # 说明。
  # 说明。
  # 说明。
  info "正在拉取新镜像（仅 artex）…"
  docker compose pull artex
  info "正在重新创建并启动（artex 会在重启时自动迁移 schema）…"
  docker compose up -d artex
  ok "更新完成 → http://localhost:8787"
  warn "刚拉取的是上游 autumn27/artex 镜像，不包含此仓库的本地化改动。若要使用当前代码，请选择“2) 本地更新（使用 go 重新编译）”"
  info "查看日志：docker compose logs -f artex"
  info "清理旧镜像（可选）：docker image prune -f"
}

# 说明。
update_local(){
  command -v go >/dev/null 2>&1 || die "未找到 Go（需要 >=1.26）：https://go.dev/dl/"
  [ -f config.json ] || warn "找不到 config.json。首次部署请运行 ./install.sh"
  ok "Go: $(go version)"

  if command -v npm >/dev/null 2>&1; then
    info "正在重新构建前端静态文件…"
    ( cd web && npm ci && npm run build:static )
    rm -rf server/webui/dist && cp -r web/out server/webui/dist
    info "正在重新编译内嵌前端的单体二进制文件…"
    CGO_ENABLED=0 go build -tags embedui -trimpath -o artex ./cmd/artex
  else
    warn "未找到 npm，仅编译不含前端的后端（需单独运行 npm run dev 启动前端）"
    CGO_ENABLED=0 go build -o artex ./cmd/artex
  fi
  ok "编译完成 → ./artex"
  warn "请重启正在运行的 artex 进程以应用更改（重启时会自动迁移 schema）"
}

echo "=============================="
echo "  ARTEX 更新"
echo "  1) Docker 更新（拉取新镜像并重新创建容器）"
echo "  2) 本地更新（使用 go 重新编译）"
echo "=============================="
case "$(ask '请选择' 1)" in
  1) sync_repo; update_docker ;;
  2) sync_repo; update_local ;;
  *) die "选项无效" ;;
esac
