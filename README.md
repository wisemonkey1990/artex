<div align="center">

# ARTEX 中文版

**由大语言模型驱动的自主渗透测试系统**（Go 后端 + Next.js 前端）

简体中文 · [한국어](README.ko.md) · [English](README.en.md)

[![license: AGPL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)

</div>

---

> **仅限授权安全测试。** ARTEX 能够自主执行侦察、工具调用和安全验证。请只在自己拥有或获得明确书面授权的隔离环境中使用。未经授权扫描、访问或利用他人系统可能违法并造成损害。使用前请阅读下方的安全范围说明及 [LICENSE](LICENSE)。

ARTEX 由多个 LLM 智能体协作完成目标拆解、工具执行和结果记录，并通过资产图谱与任务过程视图呈现发现。项目包含 Go 服务端、Next.js Web 界面和 PostgreSQL 数据存储。此仓库以简体中文作为默认界面和智能体输出语言，保留韩语界面作为可选语言。

## 快速开始

需要 Docker 与 Docker Compose。首次启动会创建 PostgreSQL 和 ARTEX 服务；访问 `http://localhost:8787` 并按页面提示设置管理员密码。

```bash
git clone https://github.com/wisemonkey1990/artex.git
cd artex
cp .env.example .env
docker compose up -d
```

请在 `.env` 中设置数据库密码，并按需配置 LLM API 密钥。若要从源码构建，请先构建静态前端，再将产物嵌入 Go 服务端：

```bash
cd web
npm ci
npm run build:static
cd ..
rm -rf server/webui/dist
mkdir -p server/webui/dist
cp -a web/out/. server/webui/dist/
CGO_ENABLED=0 go build -tags embedui -o artex ./cmd/artex
```

## 本地化说明

- Web 默认语言为简体中文，也可通过构建变量 `NEXT_PUBLIC_LOCALE=ko` 切换为韩语。
- 时间、日期和页面元信息按中文（中国）格式显示。
- 智能体面向用户的自然语言输出使用简体中文；命令、代码、URL、请求响应及证据原文保持不变。
- `README.zh.md` 保留上游中文项目文档；本页说明当前仓库的本地化版本。

## 安全使用范围

仅在明确授权且约定了目标、时间和测试范围的环境中运行。推荐使用本地隔离的故意脆弱靶场，不要将真实用户数据带入测试。测试中产生的数据和凭据应按约定妥善保管并及时清理。安全问题请参阅 [SECURITY.md](SECURITY.md)。

## 文档与许可

- [上游原版中文文档](README.zh.md)
- [英文项目说明](README.en.md)
- [韩语项目说明](README.ko.md)
- [贡献指南](CONTRIBUTING.md)
- [AGPL-3.0 许可](LICENSE)
