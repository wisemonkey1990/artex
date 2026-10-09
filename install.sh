#!/usr/bin/env bash
# ARTEX 安装脚本：① 全部使用 Docker ② 本地编译运行
set -euo pipefail
cd "$(cd "$(dirname "$0")" && pwd)"

info(){ printf '\033[36m[*]\033[0m %s\n' "$*"; }
ok(){   printf '\033[32m[+]\033[0m %s\n' "$*"; }
warn(){ printf '\033[33m[!]\033[0m %s\n' "$*"; }
die(){  printf '\033[31m[x]\033[0m %s\n' "$*" >&2; exit 1; }
ask(){  local p="$1" d="${2:-}" a; read -rp "$p${d:+ [$d]}: " a; echo "${a:-$d}"; }
rand(){ head -c 18 /dev/urandom | base64 | tr -dc 'A-Za-z0-9' | head -c 24; }

# ── 检测 Docker 环境并按需安装 ──────────────────
ensure_docker(){
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    ok "已检测到 docker 和 docker compose"; return
  fi
  warn "未找到 docker / docker compose"
  case "$(uname -s)" in
    Linux)
      if [ "$(ask '是否自动安装 Docker？(y/n)' y)" = y ]; then
        curl -fsSL https://get.docker.com | sh
        sudo usermod -aG docker "$USER" || true
        ok "Docker 安装完成（用户组更改需重新登录后才能免 sudo 生效）"
      else
        die "请手动安装 docker 后重新运行"
      fi ;;
    Darwin) die "请在 macOS 上安装 Docker Desktop：https://www.docker.com/products/docker-desktop/" ;;
    *)      die "请手动安装 docker 后重新运行" ;;
  esac
}

# ── ① 全部使用 Docker ──────────────────────────────
install_docker(){
  ensure_docker
  if [ ! -f .env ]; then
    cp .env.example .env 2>/dev/null || true
    local pw key
    pw="$(ask 'Postgres 密码（按回车随机生成）' "$(rand)")"
    key="$(ask 'ANTHROPIC_API_KEY（可留空，之后可在界面中设置）' '')"
    sed -i.bak "s|^POSTGRES_PASSWORD=.*|POSTGRES_PASSWORD=${pw}|" .env
    sed -i.bak "s|^ANTHROPIC_API_KEY=.*|ANTHROPIC_API_KEY=${key}|" .env
    rm -f .env.bak
    ok "已生成 .env 文件（POSTGRES_PASSWORD 已设置）"
  else
    info "沿用现有 .env 文件"
  fi
  info "正在拉取镜像并启动…"
  docker compose pull || true
  docker compose up -d
  ok "启动完成 → http://localhost:8787"
  warn "刚拉取的是上游 autumn27/artex 镜像，不包含此仓库的本地化改动"
  warn "若要运行当前代码，请重新运行此脚本并选择“2) 本地运行（使用 go 编译）”，或按 README 中的源码构建步骤操作"
  info "查看日志：docker compose logs -f artex"
}

# ── ② 本地编译运行 ────────────────────────────
install_local(){
  echo "数据库设置方式："
  echo "  1) 连接现有 PostgreSQL"
  echo "  2) 使用 Docker 启动 PostgreSQL（需要 docker）"
  case "$(ask '请选择' 1)" in
    2)
      ensure_docker
      local pw; pw="$(ask 'Postgres 密码（按回车随机生成）' "$(rand)")"
      docker run -d --name artex-pg -p 5432:5432 \
        -e POSTGRES_USER=artex -e POSTGRES_PASSWORD="$pw" -e POSTGRES_DB=artex \
        -v artex-pg:/var/lib/postgresql/data postgres:16-alpine
      DB_HOST=127.0.0.1 DB_PORT=5432 DB_USER=artex DB_PASS="$pw" DB_NAME=artex DB_SSL=disable ;;
    *)
      DB_HOST="$(ask '数据库地址' 127.0.0.1)"
      DB_PORT="$(ask '端口' 5432)"
      DB_USER="$(ask '账号' artex)"
      DB_PASS="$(ask '密码' '')"
      DB_NAME="$(ask '数据库名称' artex)"
      DB_SSL="$(ask 'sslmode (disable/require)' disable)" ;;
  esac

  # 说明。
  cat > config.json <<JSON
{
  "database": {
    "host": "${DB_HOST}",
    "port": ${DB_PORT},
    "user": "${DB_USER}",
    "password": "${DB_PASS}",
    "dbname": "${DB_NAME}",
    "sslmode": "${DB_SSL}"
  }
}
JSON
  ok "已生成 config.json 文件"

  # 说明。
  command -v go >/dev/null 2>&1 || die "未找到 Go。请先安装 Go（>=1.26）：https://go.dev/dl/"
  ok "Go: $(go version)"

  # 说明。
  if command -v npm >/dev/null 2>&1; then
    info "正在构建前端静态文件…"
    ( cd web && npm ci && npm run build:static )
    rm -rf server/webui/dist && mkdir -p server/webui && cp -r web/out server/webui/dist
    info "正在编译内嵌前端的单体二进制文件…"
    CGO_ENABLED=0 go build -tags embedui -trimpath -o artex ./cmd/artex
  else
    warn "未找到 npm：仅编译不含前端的后端（需单独运行 npm run dev 启动前端）"
    CGO_ENABLED=0 go build -o artex ./cmd/artex
  fi
  ok "编译完成 → ./artex"

  info "正在启动…（按 Ctrl-C 退出）"
  ./artex
}

echo "=============================="
echo "  ARTEX 安装"
echo "  1) 全部使用 Docker 部署"
echo "  2) 本地运行（使用 go 编译）"
echo "=============================="
case "$(ask '请选择' 1)" in
  1) install_docker ;;
  2) install_local ;;
  *) die "选项无效" ;;
esac
